//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package utils

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParametrizedPipeline) DeepCopyInto(out *ParametrizedPipeline) {
	*out = *in
	in.Pipeline.DeepCopyInto(&out.Pipeline)
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParametrizedPipeline.
func (in *ParametrizedPipeline) DeepCopy() *ParametrizedPipeline {
	if in == nil {
		return nil
	}
	out := new(ParametrizedPipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pipeline) DeepCopyInto(out *Pipeline) {
	*out = *in
	in.PipelineRef.DeepCopyInto(&out.PipelineRef)
	in.Timeouts.DeepCopyInto(&out.Timeouts)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pipeline.
func (in *Pipeline) DeepCopy() *Pipeline {
	if in == nil {
		return nil
	}
	out := new(Pipeline)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineRef) DeepCopyInto(out *PipelineRef) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make([]Param, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineRef.
func (in *PipelineRef) DeepCopy() *PipelineRef {
	if in == nil {
		return nil
	}
	out := new(PipelineRef)
	in.DeepCopyInto(out)
	return out
}
