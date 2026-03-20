package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	defaultDevice      = "/dev/block-device"
	defaultMapperName  = "mapper"
	defaultMapperPath  = "/dev/mapper/mapper"
	defaultMountTarget = "/mnt/storage"
	defaultPIDFile     = "/dev/shm/luks-helper.pid"
	defaultMaxWait     = 3600
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <format-disk|wait-ready|sleep> [flags]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "format-disk":
		runFormatDisk(os.Args[2:])
	case "wait-ready":
		runWaitReady(os.Args[2:])
	case "sleep":
		// Debug mode: block forever without doing anything.
		// Override the container command to ["luks-helper", "sleep"] and
		// kubectl exec in to run cryptsetup/mount operations manually.
		log.Println("sleep mode: doing nothing, exec in to debug")
		waitForSignal()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		usage()
	}
}

// ── format-disk ──────────────────────────────────────────────────────────────

type diskConfig struct {
	device      string
	mapperName  string
	mapperPath  string
	mountTarget string
	pidFile     string
}

func runFormatDisk(args []string) {
	fs := flag.NewFlagSet("format-disk", flag.ExitOnError)
	cfg := diskConfig{}
	fs.StringVar(&cfg.device, "device", defaultDevice, "block device to encrypt")
	fs.StringVar(&cfg.mapperName, "mapper-name", defaultMapperName, "LUKS mapper name")
	fs.StringVar(&cfg.mapperPath, "mapper-path", defaultMapperPath, "LUKS mapper device path")
	fs.StringVar(&cfg.mountTarget, "mount-target", defaultMountTarget, "mount point for the formatted volume")
	fs.StringVar(&cfg.pidFile, "pid-file", defaultPIDFile, "path to write this process's PID for wait-ready to discover")
	fs.Parse(args)

	log.Println("format-disk starting")
	log.Printf("config: device=%s mapper-name=%s mapper-path=%s mount-target=%s pid-file=%s",
		cfg.device, cfg.mapperName, cfg.mapperPath, cfg.mountTarget, cfg.pidFile)

	pass := os.Getenv("PASS")
	if pass == "" {
		fatal("PASS environment variable is not set")
	}

	cfg.mustLuksFormat(pass)
	cfg.mustLuksOpen(pass)
	cfg.mustFormatXFS()
	cfg.mustMount()
	cfg.mustWritePIDFile()
	cfg.mustMarkReady()

	log.Println("storage ready, keeping sidecar alive")
	waitForSignal()
}

func (c *diskConfig) mustLuksFormat(pass string) {
	log.Printf("checking LUKS header on %s", c.device)
	if runQuiet("cryptsetup", "isLuks", c.device) == nil {
		log.Println("device is LUKS already, skipping format")
		return
	}
	log.Printf("no LUKS header found, formatting %s with LUKS", c.device)
	if err := runWithStdin(pass, "cryptsetup", "luksFormat", c.device, "-"); err != nil {
		fatal("luksFormat: %v", err)
	}
	log.Println("LUKS format complete")
}

func (c *diskConfig) mustLuksOpen(pass string) {
	log.Printf("checking if mapper %q is already open", c.mapperName)
	if out, err := exec.Command("cryptsetup", "status", c.mapperName).CombinedOutput(); err == nil {
		log.Printf("mapper %q already open, skipping luksOpen:\n%s", c.mapperName, strings.TrimSpace(string(out)))
		return
	}
	log.Printf("opening LUKS device %s as mapper %q", c.device, c.mapperName)
	if err := runWithStdin(pass, "cryptsetup", "luksOpen", c.device, c.mapperName, "-"); err != nil {
		fatal("luksOpen: %v", err)
	}
	log.Printf("mapper %q opened at %s", c.mapperName, c.mapperPath)
}

func (c *diskConfig) mustFormatXFS() {
	// Intentionally probe the mapper device, not the raw block device.
	// After luksOpen the XFS filesystem lives inside the LUKS container, so
	// blkid on the raw device would return "crypto_LUKS", never "xfs".
	// The original shell script probed the raw device — that was a latent bug.
	log.Printf("probing filesystem type on %s", c.mapperPath)
	out, err := exec.Command("blkid", "-o", "value", "-s", "TYPE", c.mapperPath).Output()
	fsType := "none"
	if err != nil {
		log.Printf("blkid returned error (treating as unformatted): %v", err)
	} else {
		fsType = strings.TrimSpace(string(out))
	}
	log.Printf("detected filesystem type: %q", fsType)

	if fsType == "xfs" {
		log.Println("filesystem is already XFS, skipping mkfs")
		return
	}
	log.Printf("formatting %s with XFS", c.mapperPath)
	if err := runVisible("mkfs.xfs", c.mapperPath); err != nil {
		fatal("mkfs.xfs: %v", err)
	}
	log.Println("XFS format complete")
}

func (c *diskConfig) mustMount() {
	if err := os.MkdirAll(c.mountTarget, 0755); err != nil {
		fatal("create mount target %s: %v", c.mountTarget, err)
	}
	log.Printf("mounting %s at %s (xfs)", c.mapperPath, c.mountTarget)
	if err := syscall.Mount(c.mapperPath, c.mountTarget, "xfs", 0, ""); err != nil {
		if errors.Is(err, syscall.EBUSY) {
			log.Printf("%s already mounted, skipping", c.mountTarget)
			return
		}
		fatal("mount: %v", err)
	}
	log.Printf("mount successful: %s -> %s", c.mapperPath, c.mountTarget)
}

func (c *diskConfig) mustWritePIDFile() {
	pid := os.Getpid()
	log.Printf("writing PID %d to %s", pid, c.pidFile)
	if err := os.WriteFile(c.pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		fatal("write pid file: %v", err)
	}
}

func (c *diskConfig) mustMarkReady() {
	readyFile := filepath.Join(c.mountTarget, ".ready")
	log.Printf("creating ready marker: %s", readyFile)
	if err := os.WriteFile(readyFile, nil, 0644); err != nil {
		fatal("create ready file: %v", err)
	}
	log.Println("ready marker created")
}

// ── wait-ready ────────────────────────────────────────────────────────────────

func runWaitReady(args []string) {
	fs := flag.NewFlagSet("wait-ready", flag.ExitOnError)
	mountTarget := fs.String("mount-target", defaultMountTarget, "mount point to watch for the ready marker")
	pidFile := fs.String("pid-file", defaultPIDFile, "path to the PID file written by format-disk")
	maxWait := fs.Int("max-wait", defaultMaxWait, "maximum seconds to wait for the ready marker")
	fs.Parse(args)

	log.Println("wait-ready starting")
	log.Printf("config: mount-target=%s pid-file=%s max-wait=%ds", *mountTarget, *pidFile, *maxWait)

	start := time.Now()
	deadline := start.Add(time.Duration(*maxWait) * time.Second)
	var readyFile string

	for {
		if readyFile == "" {
			data, err := os.ReadFile(*pidFile)
			if err == nil {
				pid, err := strconv.Atoi(string(bytes.TrimSpace(data)))
				if err != nil {
					fatal("invalid PID in %s: %v", *pidFile, err)
				}
				readyFile = fmt.Sprintf("/proc/%d/root%s/.ready", pid, *mountTarget)
				log.Printf("helper PID %d, polling for ready marker: %s", pid, readyFile)
			}
		} else if _, err := os.Stat(readyFile); err == nil {
			log.Printf("ready marker found after %ds", int(time.Since(start).Seconds()))
			return
		}

		now := time.Now()
		if now.After(deadline) {
			if readyFile == "" {
				fatal("timed out after %ds waiting for PID file %s", *maxWait, *pidFile)
			}
			fatal("timed out after %ds waiting for %s", *maxWait, readyFile)
		}
		if elapsed := int(now.Sub(start).Seconds()); elapsed%30 == 0 && elapsed > 0 {
			log.Printf("waited %ds...", elapsed)
		}
		time.Sleep(time.Second)
	}
}

// ── shared helpers ────────────────────────────────────────────────────────────

// runQuiet runs a command discarding output — used for status checks.
func runQuiet(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

// runVisible runs a command with stdout/stderr inherited.
func runVisible(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runWithStdin pipes a string to stdin (used for cryptsetup password).
func runWithStdin(input, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// waitForSignal blocks until SIGTERM or SIGINT is received.
// Used by both the format-disk sidecar (to keep the mount alive) and the
// sleep subcommand (debug mode). Handles Kubernetes pod termination cleanly.
func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	log.Printf("received signal %s, exiting", sig)
}

func fatal(format string, args ...any) {
	log.Fatalf("FATAL: "+format, args...)
}
