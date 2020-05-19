package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
)

func (s *Service) SessionDeal(hash string, sessionId uint, blockNum int) error {
	c, _, err := websocket.DefaultDialer.Dial(utiles.ProviderEndPoint, nil)
	if err != nil {
		log.Error("dial websocket error", err)
		return errors.New("fail")
	}
	defer c.Close()
	eraIndex, err := substrate.GetStorageAt(c, hash, "Staking", "CurrentEra")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	currentEra := eraIndex.ToInt()
	validatorsList, err := s.GetValidatorFromSub(c, hash)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	var countNominators int
	for rankValidator, v := range validatorsList {
		ledger, err := substrate.GetStorageAt(c, hash, "Staking", "Ledger", utiles.TrimHex(v))
		if err != nil {
			continue
		}
		ledgerValue := ledger.ToStakingLedgers()
		fmt.Println("ledger", ledger)

		validatorStash := ledgerValue.Stash
		validatorSession, _ := substrate.GetStorageAt(c, hash, "Session", "NextKeyFor", utiles.TrimHex(v))
		validatorSessionValue := validatorSession.ToMapString()

		validatorPrefs, err := substrate.GetStorageAt(c, hash, "Staking", "Validators", utiles.TrimHex(validatorStash))
		if err != nil {
			continue
		}
		exposure, _ := substrate.GetStorageAt(c, hash, "Staking", "Stakers", utiles.TrimHex(validatorStash))
		exposureValue := exposure.ToExposures()

		fmt.Println("validatorSession", validatorSession)
		fmt.Println("validatorPrefs", validatorPrefs)
		fmt.Println("exposure", exposure)
		validatorPrefsValue := validatorPrefs.ToValidatorPrefsLegacy()
		unlockBytes, _ := json.Marshal(exposureValue.Others)
		_ = s.AddSessionValidator(sessionId, rankValidator, utiles.TrimHex(v), utiles.TrimHex(validatorSessionValue["col1"]),
			utiles.TrimHex(validatorStash), ledgerValue.ActiveRing, ledgerValue.TotalRing, exposureValue.Total.Sub(exposureValue.Own),
			exposureValue.Own, string(unlockBytes), len(exposureValue.Others), validatorPrefsValue.ValidatorPaymentRatio)

		for rankNominator, nominatorInfo := range exposureValue.Others {
			_ = s.AddSessionNominator(sessionId, rankValidator, rankNominator, utiles.TrimHex(nominatorInfo.Who), nominatorInfo.Value)
			countNominators += 1
		}
	}
	startBlockNum := blockNum + substrate.GetSessionLength()
	_ = s.AddSession(sessionId, startBlockNum, blockNum, currentEra, len(validatorsList), countNominators)
	return nil
}

func (s *Service) GetValidatorFromSub(c *websocket.Conn, hash string) ([]string, error) {
	validators, err := substrate.GetStorageAt(c, hash, "Session", "Validators")
	if err != nil {
		return []string{}, err
	}
	return validators.ToStringSlice(), nil
}
