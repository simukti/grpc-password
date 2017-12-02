// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

// Client sample app
package main

import (
	"context"
	"log"

	"github.com/simukti/grpc-password/proto"
	"google.golang.org/grpc"
)

var (
	pwd    = "12345678"
	pwdNew = "87654321"
)

func createPassword(c proto.PasswordClient) *proto.PasswordCreateResponse {
	r := &proto.PasswordCreateRequest{
		Type:          proto.ARGON2,
		PlainPassword: pwd,
	}

	resp, err := c.PasswordCreate(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
	return resp
}

func validatePassword(c proto.PasswordClient, id string) *proto.PasswordValidateResponse {
	r := &proto.PasswordValidateRequest{
		Id:            id,
		PlainPassword: pwd,
	}

	resp, err := c.PasswordValidate(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
	return resp
}

func validatePasswordNew(c proto.PasswordClient, id string) *proto.PasswordValidateResponse {
	r := &proto.PasswordValidateRequest{
		Id:            id,
		PlainPassword: pwdNew,
	}

	resp, err := c.PasswordValidate(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
	return resp
}

func updatePassword(c proto.PasswordClient, id string) *proto.PasswordUpdateResponse {
	r := &proto.PasswordUpdateRequest{
		Id:               id,
		OldPlainPassword: pwd,
		NewPlainPassword: pwdNew,
	}

	resp, err := c.PasswordUpdate(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
	return resp
}

func deletePassword(c proto.PasswordClient, id string) *proto.PasswordDeleteResponse {
	r := &proto.PasswordDeleteRequest{
		Id: id,
	}

	resp, err := c.PasswordDelete(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
	return resp
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:54123", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := proto.NewPasswordClient(conn)

	log.Println("######### CREATE NEW PASSWORD CREDENTIAL #########")
	p := createPassword(client)
	log.Println("######### VALIDATE CURRENT PASSWORD #########")
	validatePassword(client, p.Id)
	log.Println("######### VALIDATE NEW PASSWORD #########")
	validatePasswordNew(client, p.Id)
	log.Println("######### UPDATE TO NEW PASSWORD #########")
	updatePassword(client, p.Id)
	log.Println("######### VALIDATE OLD PASSWORD #########")
	validatePassword(client, p.Id)
	log.Println("######### VALIDATE NEW PASSWORD #########")
	validatePasswordNew(client, p.Id)
	log.Println("######### DELETE CREDENTIAL #########")
	deletePassword(client, p.Id)
	log.Println("######### VALIDATE NEW PASSWORD #########")
	validatePassword(client, p.Id)
}
