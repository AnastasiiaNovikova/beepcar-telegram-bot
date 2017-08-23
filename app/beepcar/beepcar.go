package beepcar

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jirfag/gointensive/lec3/2_http_pool/workers"
)

var wp = workers.NewPool(2)

func init() {
	wp.Run()
}

type SuggestResponse struct {
	Result struct {
		Locations []struct {
			ID int64
		}
	}
}

func LocationNameToID(ctx context.Context, name string) (int64, error) {
	var resp SuggestResponse
	d, _ := ctx.Deadline()
	apiErr, qErr := wp.AddTaskSyncTimed(func() interface{} {
		return callAPI(fmt.Sprintf("/locations/suggest?input=%s", name), &resp)
	}, d.Sub(time.Now()))

	if qErr != nil {
		return 0, fmt.Errorf("can't add task to beepcar api: %s", qErr)
	}

	if apiErr != nil {
		return 0, fmt.Errorf("api error: %s", apiErr)
	}

	if len(resp.Result.Locations) == 0 {
		return 0, fmt.Errorf("empty locations response %+v", resp)
	}

	return resp.Result.Locations[0].ID, nil
}

type SearchResponse struct {
	Result struct {
		Trips []struct {
			TripID int64
		}
	}
}

func Search(ctx context.Context, fromLocID, toLocID int64) ([]int64, error) {
	var resp SearchResponse
	d, _ := ctx.Deadline()
	apiErr, qErr := wp.AddTaskSyncTimed(func() interface{} {
		url := fmt.Sprintf("/trips/search?from_location_id=%d&to_location_id=%d",
			fromLocID, toLocID)
		return callAPI(url, &resp)
	}, d.Sub(time.Now()))

	if qErr != nil {
		return nil, fmt.Errorf("can't add task to beepcar api: %s", qErr)
	}

	if apiErr != nil {
		return nil, fmt.Errorf("api error: %s", apiErr)
	}

	ids := []int64{}
	for _, t := range resp.Result.Trips {
		ids = append(ids, t.TripID)
	}

	return ids, nil
}

func callAPI(subPath string, ret interface{}) error {
	url := fmt.Sprintf("https://api.beepcar.ru/v1%s", subPath)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't make http request to %q: %s", url, err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read http request %q response: %s", url, err)
	}

	if err = json.Unmarshal(respBody, ret); err != nil {
		return fmt.Errorf("can't unmarshal response '%s' to %+v: %s",
			string(respBody), ret, err)
	}

	return nil
}
