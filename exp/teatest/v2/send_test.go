package teatest_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/purpose168/bubbletea-cn/v2"
	"github.com/purpose168/charm-experimental-packages-cn/exp/teatest/v2"
)

// TestAppSendToOtherProgram 测试向其他程序发送消息的功能
func TestAppSendToOtherProgram(t *testing.T) {
	m1 := &connectedModel{
		name: "m1",
	}
	m2 := &connectedModel{
		name: "m2",
	}

	tm1 := teatest.NewTestModel(t, m1, teatest.WithInitialTermSize(70, 30))
	t.Cleanup(func() {
		if err := tm1.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	tm2 := teatest.NewTestModel(t, m2, teatest.WithInitialTermSize(70, 30))
	t.Cleanup(func() {
		if err := tm2.Quit(); err != nil {
			t.Fatal(err)
		}
	})
	m1.programs = append(m1.programs, tm2)
	m2.programs = append(m2.programs, tm1)

	tm1.Type("pp")
	tm2.Type("pppp")

	tm1.Type("q")
	tm2.Type("q")

	out1 := readBts(t, tm1.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	out2 := readBts(t, tm2.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))

	if string(out1) != string(out2) {
		t.Errorf("两个模型的输出应该相同，得到:\n%v\n和:\n%v\n", string(out1), string(out2))
	}

	teatest.RequireEqualOutput(t, out1)
}

// connectedModel 是一个可以与其他程序通信的模型
type connectedModel struct {
	name     string
	programs []interface{ Send(tea.Msg) }
	msgs     []string
}

// ping 是一个表示 ping 消息的类型
type ping string

func (m *connectedModel) Init() tea.Cmd {
	return nil
}

func (m *connectedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			send := ping("from " + m.name)
			m.msgs = append(m.msgs, string(send))
			for _, p := range m.programs {
				p.Send(send)
			}
			fmt.Printf("向其他程序发送 ping %q\n", send)
		case "q":
			return m, tea.Quit
		}
	case ping:
		fmt.Printf("在 %s 上收到 ping %q\n", m.name, msg)
		m.msgs = append(m.msgs, string(msg))
	}
	return m, nil
}

func (m *connectedModel) View() string {
	return "All pings:\n" + strings.Join(m.msgs, "\n")
}
