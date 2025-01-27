// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package types

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/params"
)

var (
	// DefaultMinGasMultiplier is 0.5 or 50%
	DefaultMinGasMultiplier = sdkmath.LegacyNewDecWithPrec(50, 2)
	DefaultMinGasPrice      = sdkmath.LegacyNewDec(50_000_000_000_000)
	// DefaultEnableHeight is 0 (i.e disabled)
	DefaultEnableHeight = int64(0)
	// DefaultNoBaseFee is false
	DefaultNoBaseFee        = false
	DefaultGasLimitPerBlock = sdkmath.NewInt(1_000_000_000)
	DefaultInitialBaseFee   = sdkmath.NewInt(50_000_000_000_000)
	// DefaultMaxBaseFee = 5 * 10^19
	DefaultMaxBaseFee = sdkmath.NewIntFromBigInt(new(big.Int).Mul(new(big.Int).Exp(big.NewInt(10), big.NewInt(19), nil), big.NewInt(5)))
)

// Parameter keys
var (
	ParamsKey                             = []byte("Params")
	ParamStoreKeyNoBaseFee                = []byte("NoBaseFee")
	ParamStoreKeyBaseFeeChangeDenominator = []byte("BaseFeeChangeDenominator")
	ParamStoreKeyElasticityMultiplier     = []byte("ElasticityMultiplier")
	ParamStoreKeyBaseFee                  = []byte("BaseFee")
	ParamStoreKeyEnableHeight             = []byte("EnableHeight")
	ParamStoreKeyMinGasPrice              = []byte("MinGasPrice")
	ParamStoreKeyMinGasMultiplier         = []byte("MinGasMultiplier")
	ParamStoreKeyGasLimitPerBlock         = []byte("GasLimitPerBlock")
	ParamStoreKeyMaxBaseFee               = []byte("MaxBaseFee")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyNoBaseFee, &p.NoBaseFee, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFeeChangeDenominator, &p.BaseFeeChangeDenominator, validateBaseFeeChangeDenominator),
		paramtypes.NewParamSetPair(ParamStoreKeyElasticityMultiplier, &p.ElasticityMultiplier, validateElasticityMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFee, &p.BaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableHeight, &p.EnableHeight, validateEnableHeight),
		paramtypes.NewParamSetPair(ParamStoreKeyMinGasPrice, &p.MinGasPrice, validateMinGasPrice),
		paramtypes.NewParamSetPair(ParamStoreKeyMinGasMultiplier, &p.MinGasMultiplier, validateMinGasPrice),
		paramtypes.NewParamSetPair(ParamStoreKeyGasLimitPerBlock, &p.GasLimitPerBlock, validateGasLimitPerBlock),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxBaseFee, &p.MaxBaseFee, validateGasLimitPerBlock),
	}
}

// NewParams creates a new Params instance
func NewParams(
	noBaseFee bool,
	baseFeeChangeDenom,
	elasticityMultiplier uint32,
	baseFee uint64,
	enableHeight int64,
	minGasPrice sdkmath.LegacyDec,
	minGasPriceMultiplier sdkmath.LegacyDec,
	gasLimitPerBlock sdkmath.Int,
	maxBaseFee sdkmath.Int,
) Params {
	return Params{
		NoBaseFee:                noBaseFee,
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		BaseFee:                  sdkmath.NewIntFromUint64(baseFee),
		EnableHeight:             enableHeight,
		MinGasPrice:              minGasPrice,
		MinGasMultiplier:         minGasPriceMultiplier,
		GasLimitPerBlock:         gasLimitPerBlock,
		MaxBaseFee:               maxBaseFee,
	}
}

// DefaultParams returns default evm parameters
func DefaultParams() Params {
	return Params{
		NoBaseFee:                DefaultNoBaseFee,
		BaseFeeChangeDenominator: params.BaseFeeChangeDenominator,
		ElasticityMultiplier:     params.ElasticityMultiplier,
		BaseFee:                  DefaultInitialBaseFee,
		EnableHeight:             DefaultEnableHeight,
		MinGasPrice:              DefaultMinGasPrice,
		MinGasMultiplier:         DefaultMinGasMultiplier,
		GasLimitPerBlock:         DefaultGasLimitPerBlock,
		MaxBaseFee:               DefaultMaxBaseFee,
	}

}

// DefaultEthermintTestParams returns default evm parameters for local testing
func DefaultEthermintTestParams() Params {
	return Params{
		NoBaseFee:                false,
		BaseFeeChangeDenominator: params.BaseFeeChangeDenominator,
		ElasticityMultiplier:     params.ElasticityMultiplier,
		BaseFee:                  sdkmath.NewIntFromUint64(params.InitialBaseFee),
		EnableHeight:             0,
		MinGasPrice:              sdkmath.LegacyZeroDec(),
		MinGasMultiplier:         DefaultMinGasMultiplier,
		GasLimitPerBlock:         DefaultGasLimitPerBlock,
		MaxBaseFee:               DefaultMaxBaseFee,
	}
}

// Validate performs basic validation on fee market parameters.
func (p Params) Validate() error {
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	if p.BaseFee.IsNegative() {
		return fmt.Errorf("initial base fee cannot be negative: %s", p.BaseFee)
	}

	if p.EnableHeight < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", p.EnableHeight)
	}

	if err := validateMinGasMultiplier(p.MinGasMultiplier); err != nil {
		return err
	}

	if err := validateMinGasPrice(p.MinGasPrice); err != nil {
		return err
	}

	if err := validateGasLimitPerBlock(p.GasLimitPerBlock); err != nil {
		return err
	}

	return validateMaxBaseFee(p.MaxBaseFee)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func (p *Params) IsBaseFeeEnabled(height int64) bool {
	return !p.NoBaseFee && height >= p.EnableHeight
}

func validateMinGasPrice(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", i)
	}

	return nil
}

func validateBaseFeeChangeDenominator(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	return nil
}

func validateElasticityMultiplier(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBaseFee(i interface{}) error {
	value, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("base fee cannot be negative")
	}

	return nil
}

func validateEnableHeight(i interface{}) error {
	value, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", value)
	}

	return nil
}

func validateMinGasMultiplier(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", v)
	}

	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %s", v)
	}
	return nil
}
func validateGasLimitPerBlock(i interface{}) error {
	value, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("gas limit per block cannot be negative")
	}

	return nil
}

func validateMaxBaseFee(i interface{}) error {
	value, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("max base fee cannot be negative")
	}

	return nil
}
