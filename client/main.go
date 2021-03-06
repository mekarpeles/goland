package main

import (
	"flag"
	//"github.com/aarzilli/golua/lua"
	"github.com/mischief/goland/game/gutil"
	lua "github.com/xenith-studios/golua"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	configfile = flag.String("config", "config.lua", "configuration file")

	Lua *lua.State
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Lua = lua.NewState()

	if Lua == nil {
		log.Fatal("main: init: Can't make lua state")
	}

	Lua.OpenLibs()
}

func main() {
	flag.Parse()

	// load configuration
	ParMap, err := gutil.LuaParMapFromFile(Lua, *configfile)
	if err != nil || ParMap == nil {
		log.Fatalf("main: Error loading configuration file %s: %s", *configfile, err)
	}

	lf, ok := ParMap.Get("logfile")
	if !ok {
		log.Printf("main: No logfile specified, using stdout")
	} else {
		// open log file
		f, err := os.OpenFile(lf, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Print("main: Logging started")

	// log panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("main: Recovered from %v", r)
		}
	}()

	log.Printf("main: Config loaded from %s", *configfile)

	// dump config
	it := ParMap.Iter()
	for k, v, b := it(); b != false; k, v, b = it() {
		log.Printf("main: config: %s -> '%s'", k, v)
	}

	// enable profiling
	if cpuprofile, ok := ParMap.Get("cpuprofile"); ok {
		log.Printf("main: Starting profiling in file %s", cpuprofile)
		f, err := os.OpenFile(cpuprofile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Println("Creating game instance")
	g := NewGame(ParMap)

	g.Run()

	log.Println("-- Logging ended --")
}
