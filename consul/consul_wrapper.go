package consul
import (
	"net"
	"strconv"
	"log"
	"strings"
	"os/exec"
	"os"
)

type ConsulWrapper struct {
	baseArguments string
	additionalArguments string
}

func CreateConsulWrapper() *ConsulWrapper {
	return &ConsulWrapper{"", ""}
}

func (c *ConsulWrapper) AddAdditionalArguments(args string) {
	c.additionalArguments = args
}

func (c *ConsulWrapper) Run(bin string, ip net.IPAddr, expectedNodeCnt int, nodes []string) int {
	unquotedArgs,err := strconv.Unquote(c.additionalArguments);
	if(err != nil) {
		log.Printf("Invalid arguments detected: %v\n", err)
		return 3
	}

	cli := strings.Join([]string{"agent -server ", unquotedArgs," -bootstrap-expect ", strconv.Itoa(expectedNodeCnt), " -bind ", ip.String(), " join ", strings.Join(nodes, " ")}, "")
	log.Printf("Executing child process:\n\t%s %s", bin, cli)
	cmd := exec.Command(bin, strings.Split(cli, " ")...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return 4
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
		return 5
	} else {
		return 0
	}
}