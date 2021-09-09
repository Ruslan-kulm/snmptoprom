package main

import "github.com/sirupsen/logrus"

type Context struct {
	config *Config
	logger *logrus.Logger
}
