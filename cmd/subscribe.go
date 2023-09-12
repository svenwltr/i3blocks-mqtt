package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

type SubscribeRunner struct {
	broker string
	topic  string

	json          bool
	templateFull  string
	templateShort string
	templateColor string
}

func (r *SubscribeRunner) Bind(cmd *cobra.Command) error {
	cmd.PersistentFlags().StringVar(
		&r.broker, "broker", "",
		`Specify the MQTT broker URI (e.g., tcp://10.10.1.1:1883)`)
	cmd.PersistentFlags().StringVar(
		&r.topic, "topic", "",
		`Specify the topic to subscribe to`)

	cmd.PersistentFlags().BoolVar(
		&r.json, "json", false,
		`Parse the message payload as JSON before passing to the template`)

	cmd.PersistentFlags().StringVar(
		&r.templateFull, "template-full", "{{.}}",
		`Specify the Go template to format the full_text`)
	cmd.PersistentFlags().StringVar(
		&r.templateShort, "template-short", "",
		`Specify the Go template to format the short_text`)
	cmd.PersistentFlags().StringVar(
		&r.templateColor, "template-color", "",
		`Specify the Go template to format the color`)

	cmd.PreRun = AutoEnv

	return nil
}

func (r *SubscribeRunner) Run(ctx context.Context) error {
	templateFull, err := template.New("full_text").Parse(r.templateFull)
	if err != nil {
		return fmt.Errorf("parse message_full: %w", err)
	}

	templateShort, err := template.New("short_text").Parse(r.templateShort)
	if err != nil {
		return fmt.Errorf("parse message_full: %w", err)
	}

	templateColor, err := template.New("message_full").Parse(r.templateColor)
	if err != nil {
		return fmt.Errorf("parse message_full: %w", err)
	}

	handler := SubscribeHandler{
		TemplateFull:  templateFull,
		TemplateShort: templateShort,
		TemplateColor: templateColor,
		ParseJSON:     r.json,
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(r.broker)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	opts.OnConnect = func(client mqtt.Client) {
		token := client.Subscribe(r.topic, 1, handler.Handle)
		token.Wait()
		if token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	opts.OnConnectionLost = handler.onConnectionLost

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	<-ctx.Done()
	return nil
}

type SubscribeHandler struct {
	TemplateFull  *template.Template
	TemplateShort *template.Template
	TemplateColor *template.Template
	ParseJSON     bool
}

func (h *SubscribeHandler) Handle(client mqtt.Client, message mqtt.Message) {
	err := h.handle(message.Payload())
	if err != nil {
		fmt.Println(err)
	}
}

type Line struct {
	FullText  string `json:"full_text"`
	ShortText string `json:"short_text,omitempty"`
	Color     string `json:"color,omitempty"`
}

func (h *SubscribeHandler) handle(payload []byte) error {
	var data any
	if h.ParseJSON {
		err := json.Unmarshal(payload, &data)
		if err != nil {
			return fmt.Errorf("parse json: %w", err)
		}
	} else {
		data = string(payload)
	}

	var (
		err    error
		result Line
	)

	result.FullText, err = h.executeTemplate(h.TemplateFull, data)
	if err != nil {
		return fmt.Errorf("execute full_text: %w", err)
	}

	result.ShortText, err = h.executeTemplate(h.TemplateShort, data)
	if err != nil {
		return fmt.Errorf("execute short_text: %w", err)
	}

	result.Color, err = h.executeTemplate(h.TemplateColor, data)
	if err != nil {
		return fmt.Errorf("execute color: %w", err)
	}

	out, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	fmt.Println(string(out))

	return nil
}

func (h *SubscribeHandler) executeTemplate(t *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (h *SubscribeHandler) onConnectionLost(mqtt.Client, error) {
	fmt.Println("connection lost")
}

type Logger string

func (l Logger) Println(v ...interface{}) {
	fmt.Print(l, " ")
	fmt.Println(v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	fmt.Print(l, " ")
	fmt.Printf(format, v...)
}
