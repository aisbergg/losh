package dgraph

import (
	"context"
	"losh/internal/errors"
	"losh/internal/models"
	"losh/internal/repository"
)

var (
	errGetHostStr    = "failed to get host(s)"
	errSaveHostStr   = "failed to save host(s)"
	errDeleteHostStr = "failed to delete host(s)"
)

// GetHost returns a `Host` object by its ID.
func (dr *DgraphRepository) GetHost(id, domain *string) (*models.Host, error) {
	ctx := context.Background()
	getHost, err := dr.client.GetHost(ctx, id, domain)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetHostStr).
			Add("hostId", id).Add("hostDomain", domain)
	}
	if getHost.GetHost == nil { // not found
		return nil, nil
	}
	host := &models.Host{ID: *id}
	if err = dr.dataCopier.CopyTo(getHost.GetHost, host); err != nil {
		panic(err)
	}
	return host, nil
}

// GetHosts returns a list of `Host` objects matching the filter criteria.
func (dr *DgraphRepository) GetHosts(filter *models.HostFilter, order *models.HostOrder, first *int64, offset *int64) ([]*models.Host, error) {
	ctx := context.Background()
	getHosts, err := dr.client.GetHosts(ctx, filter, order, first, offset)
	if err != nil {
		return nil, repository.WrapRepoError(err, errGetHostStr)
	}
	hosts := make([]*models.Host, 0, len(getHosts.QueryHost))
	for _, x := range getHosts.QueryHost {
		host := &models.Host{ID: x.ID}
		if err = dr.dataCopier.CopyTo(x, host); err != nil {
			panic(err)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

// GetAllHosts returns a list of all `Host` objects.
func (dr *DgraphRepository) GetAllHosts() ([]*models.Host, error) {
	return dr.GetHosts(nil, nil, nil, nil)
}

// SaveHost saves a `Host` object if does not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveHost(host *models.Host) (err error) {
	err = dr.SaveHosts([]*models.Host{host})
	if aerr, ok := err.(errors.ContextAdder); ok {
		// enrich error context
		aerr.Add("hostId", host.ID).Add("hostDomain", host.Domain)
	}
	return
}

// SaveHosts saves `Host` objects which do not have an ID. After saving, the ID
// field of the input is set to the ID of the `Database` object.
func (dr *DgraphRepository) SaveHosts(hosts []*models.Host) error {
	reqData := make([]*models.AddHostInput, 0, len(hosts))
	for _, x := range hosts {
		if x.ID != "" {
			continue
		}
		host := &models.AddHostInput{}
		if err := dr.dataCopier.CopyTo(x, host); err != nil {
			return repository.WrapRepoError(err, errSaveHostStr).
				Add("hostId", x.ID).Add("hostDomain", x.Domain)
		}
		reqData = append(reqData, host)
	}
	if len(reqData) == 0 {
		return nil
	}
	ctx := context.Background()
	respData, err := dr.client.SaveHosts(ctx, reqData)
	if err != nil {
		return repository.WrapRepoError(err, errSaveHostStr)
	}
	// save ID from response
	for i, x := range hosts {
		x.ID = respData.AddHost.Host[i].ID
	}
	return nil
}

// DeleteHost deletes a `Host` object.
func (dr *DgraphRepository) DeleteHost(id, domain *string) error {
	ctx := context.Background()
	delFilter := models.HostFilter{}
	if id != nil {
		delFilter.ID = []string{*id}
	}
	if domain != nil {
		delFilter.Domain = &models.StringHashFilterStringRegExpFilter{Eq: domain}
	}
	_, err := dr.client.DeleteHost(ctx, delFilter)
	if err != nil {
		return repository.WrapRepoError(err, errDeleteHostStr).
			Add("hostId", id).Add("hostDomain", domain)
	}
	return nil
}

// DeleteAllHosts deletes all `Hosts` objects.
func (dr *DgraphRepository) DeleteAllHosts() error {
	return dr.DeleteHost(nil, nil)
}

// saveHostIfNecessary saves a `Host` object if it is not already saved.
func (dr *DgraphRepository) saveHostIfNecessary(host *models.Host) (*models.HostRef, error) {
	if host == nil {
		return nil, nil
	}
	if host.ID == "" {
		if err := dr.SaveHost(host); err != nil {
			return nil, err
		}
	}
	return &models.HostRef{ID: &host.ID}, nil
}
