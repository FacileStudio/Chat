package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/ui"
	"github.com/FacileStudio/Chat/internal/xmpp"
	tea "github.com/charmbracelet/bubbletea"
)

// Execute runs the main command-line routine.
func Execute() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	showVersion, showHelp, err := parseFlags(cfg)
	if err != nil {
		return err
	}
	if showVersion {
		fmt.Printf("chat %s\n", Version())
		return nil
	}
	if showHelp {
		printHelp()
		return nil
	}
	return runApp(cfg)
}

func parseFlags(cfg *config.Config) (bool, bool, error) {
	fs := flag.NewFlagSet("chat", flag.ContinueOnError)
	var showVersion, showHelp bool
	var configPath string
	setupFlags(fs, cfg, &showVersion, &showHelp, &configPath)
	if err := fs.Parse(os.Args[1:]); err != nil {
		return false, false, err
	}
	if configPath != "" {
		if err := applyConfigFile(configPath, cfg); err != nil {
			return false, false, err
		}
	}
	return showVersion, showHelp, nil
}

func setupFlags(fs *flag.FlagSet, cfg *config.Config, ver, help *bool, path *string) {
	fs.StringVar(path, "config", "", "path to config file")
	fs.StringVar(path, "c", "", "path to config file")
	fs.StringVar(&cfg.JID, "jid", cfg.JID, "user jid")
	fs.StringVar(&cfg.JID, "j", cfg.JID, "user jid")
	fs.StringVar(&cfg.Password, "password", cfg.Password, "user password")
	fs.StringVar(&cfg.Password, "p", cfg.Password, "user password")
	fs.StringVar(&cfg.Server, "server", cfg.Server, "server host")
	fs.StringVar(&cfg.Server, "s", cfg.Server, "server host")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "server port")
	fs.BoolVar(&cfg.InsecureTLS, "insecure", cfg.InsecureTLS, "allow insecure tls")
	fs.BoolVar(ver, "version", false, "show version")
	fs.BoolVar(ver, "v", false, "show version")
	fs.BoolVar(help, "help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
}

func applyConfigFile(path string, cfg *config.Config) error {
	loaded, err := config.LoadFile(path)
	if err != nil {
		return err
	}
	*cfg = *loaded
	return nil
}

func printHelp() {
	fmt.Printf("chat %s - Facile Studio XMPP terminal chat client\n\n", Version())
	fmt.Println("Usage:")
	fmt.Println("  chat [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -j, --jid <jid>         Account JID (e.g. user@example.com)")
	fmt.Println("  -p, --password <pass>   Account password")
	fmt.Println("  -s, --server <host>     Server hostname (optional, defaults to JID domain)")
	fmt.Println("      --port <port>       Server port (default: 5222)")
	fmt.Println("      --insecure          Skip TLS certificate verification")
	fmt.Println("  -c, --config <file>     Path to YAML config file")
	fmt.Println("  -v, --version           Print version and exit")
	fmt.Println("  -h, --help              Show help information")
}

func runApp(cfg *config.Config) error {
	events := make(chan any, 100)
	var client *xmpp.Client
	if cfg.JID != "" {
		var err error
		client, err = xmpp.New(cfg, events)
		if err != nil {
			return err
		}
		defer client.Close()
	}
	model := ui.New(cfg, client, events)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
