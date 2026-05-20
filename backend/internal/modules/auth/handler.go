package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	platformvalidator "authapp/internal/platform/validator"
)

// Handler binds HTTP requests to auth service.
type Handler struct {
	svc *Service
	val *platformvalidator.V
}

func NewHandler(svc *Service, val *platformvalidator.V) *Handler {
	return &Handler{svc: svc, val: val}
}

// SignUp godoc
// @Summary      Sign up
// @Description  Register a new user and receive a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      SignUpRequest true  "Signup payload"
// @Success      201   {object}  AuthTokensResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if errs := h.val.ValidateStruct(&req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}
	out, err := h.svc.SignUp(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		st := signupHTTPStatus(err)
		if st >= http.StatusInternalServerError {
			c.JSON(st, gin.H{"error": "internal error"})
			return
		}
		c.JSON(st, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// SignIn godoc
// @Summary      Sign in
// @Description  Authenticate and receive a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      SignInRequest true "Credentials"
// @Success      200   {object}  AuthTokensResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/signin [post]
func (h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	if errs := h.val.ValidateStruct(&req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}
	out, err := h.svc.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := signinHTTPStatus(err)
		if status == http.StatusUnauthorized {
			c.JSON(status, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(status, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, out)
}
