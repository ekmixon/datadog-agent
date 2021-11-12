// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.Datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tagset

import "sync"

// threadsafeFactory wraps another factory and uses a mutex to control
// access.
type threadsafeFactory struct {
	sync.Mutex
	Factory
}

// NewThreadsafeFactory wraps the given factory with a mutex, ensuring
// thread-safe operation.
func NewThreadsafeFactory(inner Factory) Factory {
	return &threadsafeFactory{
		Factory: inner,
	}
}

// NewTags implements Factory.NewTags
func (f *threadsafeFactory) NewTags(src []string) *Tags {
	f.Lock()
	tags := f.Factory.NewTags(src)
	f.Unlock()
	return tags
}

// NewUniqueTags implements Factory.NewUniqueTags
func (f *threadsafeFactory) NewUniqueTags(src ...string) *Tags {
	f.Lock()
	tags := f.Factory.NewUniqueTags(src...)
	f.Unlock()
	return tags
}

// NewTagsFromMap implements Factory.NewTagsFromMap
func (f *threadsafeFactory) NewTagsFromMap(src map[string]struct{}) *Tags {
	f.Lock()
	tags := f.Factory.NewTagsFromMap(src)
	f.Unlock()
	return tags
}

// NewTag implements Factory.NewTag
func (f *threadsafeFactory) NewTag(tag string) *Tags {
	f.Lock()
	tags := f.Factory.NewTag(tag)
	f.Unlock()
	return tags
}

// NewBuilder implements Factory.NewBuilder
func (f *threadsafeFactory) NewBuilder(capacity int) *Builder {
	f.Lock()
	tags := f.Factory.NewBuilder(capacity)
	f.Unlock()
	return tags
}

// ParseDSD implements Factory.ParseDSD
func (f *threadsafeFactory) ParseDSD(data []byte) (*Tags, error) {
	f.Lock()
	tags, err := f.Factory.ParseDSD(data)
	f.Unlock()
	return tags, err
}

// Union implements Factory.Union
func (f *threadsafeFactory) Union(a, b *Tags) *Tags {
	f.Lock()
	tags := f.Factory.Union(a, b)
	f.Unlock()
	return tags
}

// DisjointUnion implements Factory.DisjoingUnion
func (f *threadsafeFactory) DisjointUnion(a, b *Tags) *Tags {
	f.Lock()
	tags := f.Factory.DisjointUnion(a, b)
	f.Unlock()
	return tags
}

// getCachedTags implements Factory.getCachedTags
func (f *threadsafeFactory) getCachedTags(cacheID cacheID, hash uint64, miss func() *Tags) *Tags {
	f.Lock()
	tags := f.Factory.getCachedTags(cacheID, hash, miss)
	f.Unlock()
	return tags
}

// builderClosed implements Factory.builderClosed
func (f *threadsafeFactory) builderClosed(builder *Builder) {
	f.Lock()
	f.Factory.builderClosed(builder)
	f.Unlock()
}
