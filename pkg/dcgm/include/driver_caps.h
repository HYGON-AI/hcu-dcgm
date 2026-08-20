/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */

#ifndef DCGM_DRIVER_CAPS_H
#define DCGM_DRIVER_CAPS_H
#include <dlfcn.h>

/* 仅在 librocm_smi64.so 内查符号；勿用 RTLD_DEFAULT，否则会把
 * 可执行文件里未解析的 PLT 条目误判为“存在”（链接用了
 * --unresolved-symbols=ignore-in-object-files）。 */
static void *dcgm_rsmi_handle(void) {
	static void *handle = (void *)-1;
	if (handle == (void *)-1) {
		handle = dlopen("librocm_smi64.so", RTLD_LAZY | RTLD_NOLOAD);
		if (handle == NULL) {
			handle = dlopen("/opt/hyhal/lib/librocm_smi64.so", RTLD_LAZY | RTLD_NOLOAD);
		}
		if (handle == NULL) {
			handle = dlopen("librocm_smi64.so", RTLD_LAZY);
		}
		if (handle == NULL) {
			handle = dlopen("/opt/hyhal/lib/librocm_smi64.so", RTLD_LAZY);
		}
		if (handle == NULL) {
			handle = NULL;
		}
	}
	return handle;
}

static inline void *dcgm_lookup_symbol(const char *name) {
	void *handle = dcgm_rsmi_handle();
	void *sym;

	if (handle == NULL) {
		return NULL;
	}
	dlerror();
	sym = dlsym(handle, name);
	if (dlerror() != NULL) {
		return NULL;
	}
	return sym;
}

#endif /* DCGM_DRIVER_CAPS_H */
