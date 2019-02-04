package controller

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	logger    *zap.Logger
	clientSet kubernetes.Interface
}

func New(logger *zap.Logger, clientSet kubernetes.Interface) (*Controller, error) {
	return &Controller{
		logger:    logger,
		clientSet: clientSet,
	}, nil
}

// Run starts the controller
func (c *Controller) Run(stopChan <-chan struct{}) {
	c.logger.Info("starting controller...")

	for {
		c.logger.Debug("running control loop...")
		if err := c.runOnce(); err != nil {
			c.logger.Error("failed to runOnce", zap.Error(err))
		}

		select {
		case <-time.After(10 * time.Second):
		case <-stopChan:
			c.logger.Info("terminating controller...")
			return
		}
	}
}

func (c *Controller) runOnce() error {
	nss, err := c.clientSet.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "faild to get Namespaces")
	}
	for _, ns := range nss.Items {
		c.logger.Info(ns.Name)

		pdbs, err := c.clientSet.PolicyV1beta1().PodDisruptionBudgets(ns.Name).List(v1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "faild to get PodDistruptionBudgets in namespace: %s", ns.Name)
		}
		if len(pdbs.Items) == 0 {
			c.logger.Info(fmt.Sprintf("no PodDistruptionBudget in namespace: %s", ns.Name))
			continue
		}
		for _, pdb := range pdbs.Items {
			c.logger.Info(pdb.Name)
		}

		// TODO...
	}
	return nil
}
