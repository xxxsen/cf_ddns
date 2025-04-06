package notifier

import (
	"bytes"
	"cf_ddns/model"
	"cf_ddns/notifier"
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/xxxsen/common/utils"
)

type tgMsgNotifier struct {
	c    *config
	tplt *template.Template
}

func NewTGMsgNotifier(opts ...Option) (notifier.INotifier, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.user) == 0 || len(c.code) == 0 || len(c.addr) == 0 {
		return nil, fmt.Errorf("invalid params")
	}
	tplt := template.New("notification").Funcs(map[string]interface{}{
		"TsPrint": tsPrint,
	})
	executor, err := tplt.Parse(htmlTplt)
	if err != nil {
		return nil, err
	}

	return &tgMsgNotifier{c: c, tplt: executor}, nil
}

func (n *tgMsgNotifier) Name() string {
	return notifier.NameTGMsg
}

func (n *tgMsgNotifier) Notify(ctx context.Context, nt *model.Notification) error {
	msg, err := renderMsg(n.tplt, nt)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.c.addr, bytes.NewReader([]byte(msg)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(n.c.user, n.c.code)
	req.Header.Add("x-message-type", "tg_raw_html")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code:%d", rsp.StatusCode)
	}
	return nil
}

func createNotifier(args interface{}) (notifier.INotifier, error) {
	c := &jsConfig{}
	if err := utils.ConvStructJson(args, c); err != nil {
		return nil, err
	}
	return NewTGMsgNotifier(WithAddr(c.Addr), WithAuth(c.User, c.Code))
}

func init() {
	notifier.Register(notifier.NameTGMsg, createNotifier)
}
