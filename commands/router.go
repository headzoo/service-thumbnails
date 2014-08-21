package commands

import (
	"errors"
	"fmt"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"strings"
	"path/filepath"
)

// Commanders is a map of Commander instances.
type Commanders map[string]Commander

// Router is used to dispatch command line instructions to executors.
type Router struct {
	coms    Commanders
	inFiles []string
	outFile string
}

// NewRouter creates and returns a new Router instance.
func NewRouter(inFiles []string, outFile string) *Router {
	return &Router{
		coms:    make(Commanders),
		inFiles: inFiles,
		outFile: outFile,
	}
}

// Command registers a command executor with the router.
func (r *Router) Command(ins string, exec Commander) {
	r.coms[ins] = exec
}

// Execute executes the given instruction for the given in and out files.
func (r *Router) Route(ins string) error {
	cmd, ok := r.coms[ins]
	if !ok {
		return errors.New("No command executor for instruction " + ins)
	}

	cf := make(ChannelFinished)
	ce := make(ChannelError)
	cmd.SetChannels(&cf, &ce)

	thumbnailer.Verbose("Generating %d thumbnail(s).", len(r.inFiles))
	for i, fin := range r.inFiles {
		base := strings.TrimSuffix(fin, filepath.Ext(fin))
		fout := expandFileName(r.outFile, base, ins, i)
		go cmd.Execute(fin, fout)
	}

	running := len(r.inFiles)
	for {
		select {
		case err := <-ce:
			{
				return err
			}
		case <-cf:
			{
				running--
				if running == 0 {
					return nil
				}
			}
		}
	}

	return nil
}

// expandFileName transforms a format into a file name.
// The format may use %d which is replaced by 'index'. It may also have
// {name} which is replaced by 'name'. It may also have {type} which
// is replaced by typ.
func expandFileName(format, name, typ string, index int) string {
	if strings.Contains(format, "%") {
		format = fmt.Sprintf(format, index)
	}
	format = strings.Replace(format, "{name}", name, -1)
	format = strings.Replace(format, "{type}", typ, -1)
	
	return format
}
