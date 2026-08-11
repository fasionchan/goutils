package httpx

import (
	"context"
	"fmt"

	"github.com/fasionchan/goutils/libs/datarender"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
)

type GenericHttpTemplate[ResponseData datarender.GenericTemplatedData[DataRender_, Data], DataRender_ datarender.DataRender[Data], Data any] struct {
	RequestTemplate  *RequestTemplate `bson:"RequestTemplate" json:"RequestTemplate"`
	ResponseTemplate ResponseData     `bson:"ResponseTemplate" json:"ResponseTemplate"`
}

func (tpl *GenericHttpTemplate[ResponseData, DataRender_, Data]) Do(ctx context.Context, client *Client, funcMap templatex.TemplateFuncMap, data types.SmartJsonMap) (rawResult any, renderedResult Data, err error) {
	if tpl == nil {
		return
	}

	request, err := client.ParseAndBuildRequestWithTemplate(ctx, tpl.RequestTemplate, funcMap, data)
	if err != nil {
		err = fmt.Errorf("ParseAndBuildRequestWithTemplateFailed: %w", err)
		return
	}

	resp, err := client.DoRequestForResult(request, &rawResult, ContentTypeApplicationJson)
	if err != nil {
		err = fmt.Errorf("DoRequestForResultFailed: %w", err)
		return
	}

	rawResult = newHttpResponseDigestForRaw(resp).WithBody(rawResult)

	if data == nil {
		data = types.NewSmartJsonMap()
	}

	data.With("ResponseResult", rawResult)

	renderedResult, err = tpl.ResponseTemplate.ParseAndRender(funcMap, "", data)
	return
}
