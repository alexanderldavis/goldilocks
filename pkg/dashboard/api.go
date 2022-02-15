package dashboard

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"

	"k8s.io/klog"
	"github.com/fairwindsops/goldilocks/pkg/summary"
)

func Api(opts Options) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var namespace string
		if val, ok := vars["namespace"]; ok {
			namespace = val
		}

		summarizer := summary.NewSummarizer(
			summary.ForNamespace(namespace),
			summary.ForVPAsWithLabels(opts.vpaLabels),
			summary.ExcludeContainers(opts.excludedContainers),
		)

		vpaData, err := summarizer.GetSummary()
		if err != nil {
			klog.Errorf("Error getting vpaData: %v", err)
			http.Error(w, "Error running summary.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vpaData)
	})
}