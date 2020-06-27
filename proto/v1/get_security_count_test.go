package v1

import (
	"testing"
)

func TestNewSecurityCountRequest(t *testing.T) {
	t.Log(NewGetSecurityCount(MarketShenZhen))
	t.Log(NewGetSecurityCount(MarketShangHai))
}
