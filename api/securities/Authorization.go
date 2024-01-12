package securities

import (
	"after-sales/api/payloads"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractAuthToken(c *gin.Context) (*payloads.UserDetail, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId := fmt.Sprintf("%v", claims["user_id"])
		username := fmt.Sprintf("%s", claims["username"])
		authorized := fmt.Sprintf("%s", claims["authorized"])
		role := fmt.Sprintf("%v", claims["role"])
		companyCode := fmt.Sprintf("%s", claims["company_code"])
		ipAddress := fmt.Sprintf("%s", claims["ip_address"])
		client := fmt.Sprintf("%s", claims["client"])

		roles, _ := strconv.Atoi(role)
		userIDs, _ := strconv.Atoi(userId)
		if err != nil {
			return nil, err
		}

		authDetail := payloads.UserDetail{
			UserId:      int32(userIDs),
			Username:    username,
			Authorize:   authorized,
			Role:        uint16(roles),
			CompanyCode: companyCode,
			Client:      client,
			IpAddress:   ipAddress,
		}

		return &authDetail, nil
	}

	return nil, err
}
