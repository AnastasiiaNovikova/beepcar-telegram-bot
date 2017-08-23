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

func callAPI(subPath string, ret interface{}) error {
	url := fmt.Sprintf("https://beepcar.ru/v1%s", subPath)
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
