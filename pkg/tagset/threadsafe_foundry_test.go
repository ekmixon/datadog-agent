// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.Datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tagset

import (
	"testing"
)

func TestThreadsafeFactory(t *testing.T) {
	testFactory(t, func() Factory { return NewThreadsafeFactory(newCachingFactory()) })
	testFactoryCaching(t, func() Factory { return NewThreadsafeFactory(newCachingFactory()) })
}
