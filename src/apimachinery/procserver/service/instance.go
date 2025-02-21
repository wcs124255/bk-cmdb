package service

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateServiceInstance create service instances
func (s *service) CreateServiceInstance(ctx context.Context, h http.Header, data *metadata.CreateServiceInstanceInput) (
	[]int64, errors.CCErrorCoder) {

	resp := new(metadata.CreateServiceInstanceResp)
	subPath := "/create/proc/service_instance"

	err := s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return nil, resp.CCError()
	}
	return resp.ServiceInstanceIDs, nil
}

func (s *service) UpdateServiceInstances(ctx context.Context, h http.Header, bizID int64, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/updatemany/proc/service_instance/biz/%d"

	err = s.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) DeleteServiceInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/deletemany/proc/service_instance"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchServiceInstance search service instances
func (s *service) SearchServiceInstance(ctx context.Context, h http.Header,
	data *metadata.GetServiceInstanceInModuleInput) (*metadata.MultipleServiceInstance, error) {

	resp := new(metadata.MultipleServiceInstanceResult)
	subPath := "/findmany/proc/service_instance"

	err := s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return nil, resp.CCError()
	}
	return &resp.Data, nil
}

func (s *service) SearchServiceInstanceBySetTemplate(ctx context.Context, appID string, h http.Header, data map[string]interface{}) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)
	subPath := "/findmany/proc/service/set_template/list_service_instance/biz/%s"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) DiffServiceInstanceWithTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/proc/service_instance/difference"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) SyncServiceInstanceByTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/proc/service_instance/sync"

	err = s.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) ServiceInstanceAddLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/createmany/proc/service_instance/labels"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) ServiceInstanceRemoveLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/deletemany/proc/service_instance/labels"

	err = s.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (s *service) ServiceInstanceFindLabels(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/findmany/proc/service_instance/labels/aggregation"

	err = s.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
