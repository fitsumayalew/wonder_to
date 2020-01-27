package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"net/http"
	"xCut/util"

	. "xCut/constants"
	"xCut/entity"
	"xCut/form"
	"xCut/permission"
	"xCut/rtoken"
	"xCut/session"
	"xCut/user"
)

// UserHandler handler handles user related requests
type UserHandler struct {
	tmpl           *template.Template
	userService    user.UserService
	sessionService user.SessionService
	roleService    user.RoleService
	csrfSignKey    []byte
}


const ctxUserSessionKey = "signed_in_user_session"

type contextKey string

// NewUserHandler returns new UserHandler object
func NewUserHandler(
	t *template.Template,
	userService user.UserService,
	sessionService user.SessionService,
	roleService user.RoleService,
	csKey []byte,
) *UserHandler {
	return &UserHandler{
		tmpl:           t,
		userService:    userService,
		sessionService: sessionService,
		roleService:    roleService,
		csrfSignKey:    csKey,
	}
}

// Authenticated checks if a user is authenticated to access a given route
func (userHandler *UserHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		activeSession := userHandler.IsLoggedIn(r)
		if activeSession == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, activeSession)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (userHandler *UserHandler) getSigningKey(token *jwt.Token) (interface{}, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sessionId := claims["sessionId"].(string)
		session, err := userHandler.sessionService.Session(sessionId)
		if len(err) > 0 {
			return nil, err[0]
		}
		return session.SigningKey, nil
	}
	return nil, nil
}

func (userHandler *UserHandler) IsLoggedIn(r *http.Request) *entity.Session {
	signedStringCookie, err := r.Cookie(session.SessionKey)
	if err != nil {
		return nil
	}

	sessionId := rtoken.GetSessionIdFromToken(signedStringCookie.Value, userHandler.getSigningKey)
	if sessionId == "" {
		return nil
	}

	activeSession, errs := userHandler.sessionService.Session(sessionId)
	if len(errs) > 0 {
		return nil
	}

	return activeSession
}

func (userHandler *UserHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		///Get the user for the current active session
		activeSession := r.Context().Value(ctxUserSessionKey).(*entity.Session)
		user, errs := userHandler.userService.User(activeSession.UUID)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		///Get the role of the user
		role, errs := userHandler.roleService.Role(user.RoleID)

		//Check if the user role is authorized to access the specific path and method requested
		if len(errs) > 0 || permission.HasPermission(role.Name, r.URL.Path, r.Method) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		//Check the validity of signed token inside the form if the form is post
		if r.Method == http.MethodPost {
			if !rtoken.IsCSRFValid(r.FormValue(CsrfKey), userHandler.csrfSignKey) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	//If it's requesting the login page return CSFR Signed token with the form
	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(userHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		userHandler.tmpl.ExecuteTemplate(w, "login.layout", form.Input{
			CSRF: CSFRToken,
		})
		return
	}
	//Only reply to forms that have that are parsable and have valid csfrToken
	if util.IsParsableFormPost(w, r,userHandler.csrfSignKey) {

		//Validate form data
		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		loginForm.ValidateRequiredFields(EmailKey, PasswordKey)
		email := r.FormValue(EmailKey)
		password := r.FormValue(PasswordKey)
		user, errs := userHandler.userService.UserByEmail(email)

		///Check form validity and user password
		if len(errs) > 0 || !util.ArePasswordsSame(user.Password, password) {
			loginForm.VErrors.Add("generic", "Your email address or password is incorrect")
			userHandler.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
			return
		}

		//At this point user is successfully logged in so creating a session
		newSession, errs := userHandler.sessionService.StoreSession(session.CreateNewSession(user.ID))
		claims := rtoken.NewClaims(newSession.SessionId, newSession.Expires)
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to create session")
			userHandler.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
			return
		}
		//Save session Id in cookies
		session.SetCookie(claims, newSession.Expires, newSession.SigningKey, w)

		//Finally open the home page for the user
		if userHandler.checkAdmin(user.RoleID) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (uh *UserHandler) checkAdmin(roleID uint) bool {
	if roleID == 2 {
		return true
	}
	return false
}

// Logout logout requests
func (userHandler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//Remove cookies
	session.RemoveCookie(w)
	//Delete session from the database
	currentSession, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	userHandler.sessionService.DeleteSession(currentSession.SessionId)
	//Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func (userHandler *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CSFRToken, err := rtoken.GenerateCSRFToken(userHandler.csrfSignKey)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		userHandler.tmpl.ExecuteTemplate(w, "signup.layout", form.Input{CSRF: CSFRToken})

		return
	}
	//Only reply to forms that have that are parsable and have valid csfrToken
	if util.IsParsableFormPost(w, r,userHandler.csrfSignKey) {
		///Validate the form data
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.ValidateRequiredFields(FullnameKey, EmailKey, PasswordKey)
		signUpForm.MatchesPattern(EmailKey, form.EmailRX)
		signUpForm.MinLength(PasswordKey, 8)
		signUpForm.PasswordMatches(PasswordKey, ConfirmPasswordKey)
		signUpForm.MatchesPattern(PhoneKey, form.PhoneRX)

		if !signUpForm.IsValid() {
			userHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}
		if userHandler.userService.EmailExists(r.FormValue(EmailKey)) {
			signUpForm.VErrors.Add(EmailKey, "This email is already in use!")
			userHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}
		//Create password hash
		hashedPassword, err := util.HashPassword(r.FormValue(PasswordKey))
		if err != nil {
			signUpForm.VErrors.Add("password", "Password Could not be stored")
			userHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}
		//Create a user role for the User
		role, errs := userHandler.roleService.RoleByName("USER")
		if r.FormValue(TypeKey) == "barbershop" {
			role, errs = userHandler.roleService.RoleByName("ADMIN")
		} else {
			role, errs = userHandler.roleService.RoleByName("USER")
		}

		if len(errs) > 0 {
			signUpForm.VErrors.Add("generic", "Role couldn't be assigned to user")
			userHandler.tmpl.ExecuteTemplate(w, "signup.layout", signUpForm)
			return
		}
		///Get the data from the form and construct user object
		user := entity.User{
			FullName: r.FormValue(FullnameKey),
			Email:    r.FormValue(EmailKey),
			Phone:    r.FormValue(PhoneKey),
			Password: string(hashedPassword),
			RoleID:   role.ID,
		}
		// Save the user to the database
		_, ers := userHandler.userService.StoreUser(&user)
		if len(ers) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	}
}

