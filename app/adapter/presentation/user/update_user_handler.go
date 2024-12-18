package user

import (
	"encoding/json"
	"net/http"

	"github.com/kakkky/app/adapter/presentation/presenter"
	"github.com/kakkky/app/application/usecase/user"
	"github.com/kakkky/app/domain/errors"
	"github.com/kakkky/pkg/validation"
)

type UpdateUserHandler struct {
	updateUserUsecase *user.UpdateProfileUsecase
}

func NewUpdateUserHandler(updateUserUsecase *user.UpdateProfileUsecase) *UpdateUserHandler {
	return &UpdateUserHandler{
		updateUserUsecase: updateUserUsecase,
	}
}

// @Summary     ユーザーの更新
// @Description ユーザー情報（名前・メールアドレス）を更新する
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       request body     UpdateUserRequest                             true "ユーザー更新のための情報"
// @Success     200     {object} presenter.SuccessResponse[UpdateUserResponse] "登録されたユーザーの情報"
// @Failure     400     {object} presenter.FailureResponse                      "不正なリクエスト"
// @Failure     500     {object} presenter.FailureResponse                     "内部サーバーエラー"
// @Router      /user [patch]
func (uuh *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		presenter.RespondBadRequest(w, err.Error())
		return
	}
	if err := validation.NewValidation().Struct(&params); err != nil {
		presenter.RespondBadRequest(w, err.Error())
		return
	}
	input := user.UpdateProfileUsecaseInputDTO{
		ID:    params.ID,
		Email: params.Email,
		Name:  params.Name,
	}
	ctx := r.Context()
	output, err := uuh.updateUserUsecase.Run(ctx, input)
	if (err != nil) && errors.IsDomainErr(err) {
		presenter.RespondBadRequest(w, err.Error())
		return
	}
	if err != nil {
		presenter.RespondInternalServerError(w, err.Error())
		return
	}
	resp := UpdateUserResponse{
		ID:    output.ID,
		Email: output.Email,
		Name:  output.Name,
	}
	presenter.RespondOK(w, resp)
}