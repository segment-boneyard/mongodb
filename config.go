package main

type Config struct {
	Init         bool
	Hostname     string `valid:"host"`
	Port         string `valid:"port"`
	Username     string
	Password     string
	Database     string
}
