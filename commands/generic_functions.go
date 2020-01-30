package commands

import (
	"github.com/abuan/gitus/db"
	"github.com/spf13/cobra"
)

// Fonction à executer avant une commande pour se connecter à la BDD
func connexionForData(cmd *cobra.Command, args []string) error {
	err := db.InitDB()
	return err
}

// Fonction à executer après une commande pour se déconnecter à la BDD
func deconnexionForData(cmd *cobra.Command, args []string) error {
	db.CloseDB()
	return nil
}
