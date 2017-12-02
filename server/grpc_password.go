// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/rs/xid"
	"github.com/simukti/grpc-password/hash"
	"github.com/simukti/grpc-password/proto"
	"github.com/simukti/grpc-password/repository"
	"github.com/simukti/grpc-password/storage"
	"github.com/simukti/grpc-password/util"
)

var (
	secretStringLength = 96
	minPasswordLength  = 6
	maxPasswordLength  = 254
)

// passwordServer password rpc object
type passwordServer struct {
	db storage.Storage
}

// PasswordService password rpc request handler
func (c *grpcServer) PasswordServer() *passwordServer {
	return &passwordServer{
		db: c.db,
	}
}

func (c *passwordServer) passwordHash(hashType proto.PasswordHashType, plainPassword, secretString string) hash.PasswordHash {
	switch hashType.String() {
	case hash.Argon2Name:
		return hash.Argon2(plainPassword, secretString)
	case hash.BcryptName:
		return hash.Bcrypt(plainPassword, secretString)
	default:
		log.Fatalln(fmt.Sprintf("HASH_TYPE_NOT_AVAILABLE:%s", hashType.String()))
	}

	return nil
}

// PasswordCreate handle rpc create password
func (c *passwordServer) PasswordCreate(ctx context.Context, r *proto.PasswordCreateRequest) (*proto.PasswordCreateResponse, error) {
	strLen := utf8.RuneCountInString(r.PlainPassword)
	if strLen < minPasswordLength {
		return nil, errors.New("ERROR_TOO_SHORT")
	} else if strLen > maxPasswordLength {
		return nil, errors.New("ERROR_TOO_LONG")
	}

	secret := util.RandomString(secretStringLength)
	hasher := c.passwordHash(r.Type, r.PlainPassword, secret)
	data := &repository.PasswordModel{
		ID:        xid.New().String(),
		Type:      hasher.Name(),
		Hash:      hasher.Value(),
		Secret:    secret,
		CreatedAt: time.Now(),
	}

	err := c.db.Password().Create(data)
	if err != nil {
		return nil, err
	}

	resp := &proto.PasswordCreateResponse{Success: true, Id: data.ID}

	return resp, nil
}

// PasswordUpdate handle rpc update password
func (c *passwordServer) PasswordUpdate(ctx context.Context, r *proto.PasswordUpdateRequest) (*proto.PasswordUpdateResponse, error) {
	if _, err := xid.FromString(r.Id); err != nil {
		return nil, errors.New("INVALID_ID")
	}

	strLen := utf8.RuneCountInString(r.OldPlainPassword)
	strLen2 := utf8.RuneCountInString(r.NewPlainPassword)
	if strLen < minPasswordLength || strLen2 < minPasswordLength {
		return nil, errors.New("ERROR_TOO_SHORT")
	} else if strLen > maxPasswordLength || strLen2 > maxPasswordLength {
		return nil, errors.New("ERROR_TOO_LONG")
	}

	data, err := c.db.Password().FindByID(r.Id)
	if err != nil {
		return nil, err
	}

	if _, ok := proto.PasswordHashType_value[data.Type]; !ok {
		return nil, errors.New(fmt.Sprintf("INVALID_HASH_TYPE: %s", data.Type))
	}

	hashType := proto.PasswordHashType(proto.PasswordHashType_value[data.Type])
	hasher := c.passwordHash(hashType, r.OldPlainPassword, data.Secret)
	if !hasher.IsValid(data.Hash) {
		return nil, errors.New("ERROR_WRONG_OLD_PASSWORD")
	}

	hasher = c.passwordHash(hashType, r.NewPlainPassword, data.Secret)
	if hasher.IsValid(data.Hash) {
		return nil, errors.New("ERROR_SAME_PASSWORD")
	}

	data.Hash = hasher.Value()
	c.db.Password().Update(data)

	resp := &proto.PasswordUpdateResponse{Success: true, Id: data.ID}

	return resp, nil
}

// PasswordDelete handle rpc delete password
func (c *passwordServer) PasswordDelete(ctx context.Context, r *proto.PasswordDeleteRequest) (*proto.PasswordDeleteResponse, error) {
	if _, err := xid.FromString(r.Id); err != nil {
		return nil, errors.New("INVALID_ID")
	}

	data, err := c.db.Password().FindByID(r.Id)
	if err != nil {
		return nil, err
	}

	if err := c.db.Password().Delete(data); err != nil {
		return nil, err
	}

	resp := &proto.PasswordDeleteResponse{Success: true}

	return resp, nil
}

// PasswordValidate handle rpc password validation
func (c *passwordServer) PasswordValidate(ctx context.Context, r *proto.PasswordValidateRequest) (*proto.PasswordValidateResponse, error) {
	if _, err := xid.FromString(r.Id); err != nil {
		return nil, errors.New("INVALID_ID")
	}

	strLen := utf8.RuneCountInString(r.PlainPassword)
	if strLen < minPasswordLength {
		return nil, errors.New("ERROR_TOO_SHORT")
	} else if strLen > maxPasswordLength {
		return nil, errors.New("ERROR_TOO_LONG")
	}

	data, err := c.db.Password().FindByID(r.Id)
	if err != nil {
		return nil, err
	}

	if _, ok := proto.PasswordHashType_value[data.Type]; !ok {
		return nil, errors.New(fmt.Sprintf("INVALID_HASH_TYPE: %s", data.Type))
	}

	hashType := proto.PasswordHashType(proto.PasswordHashType_value[data.Type])
	hasher := c.passwordHash(hashType, r.PlainPassword, data.Secret)
	if !hasher.IsValid(data.Hash) {
		return nil, errors.New("ERROR_WRONG_PASSWORD")
	}

	resp := &proto.PasswordValidateResponse{Success: true, Id: data.ID}

	return resp, nil
}
