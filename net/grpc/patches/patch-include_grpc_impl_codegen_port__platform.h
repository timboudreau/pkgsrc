$NetBSD: patch-include_grpc_impl_codegen_port__platform.h,v 1.3 2019/04/19 07:15:36 adam Exp $

Add NetBSD support.

--- include/grpc/impl/codegen/port_platform.h.orig	2019-04-15 22:38:24.000000000 +0000
+++ include/grpc/impl/codegen/port_platform.h
@@ -331,6 +331,29 @@
 #else /* _LP64 */
 #define GPR_ARCH_32 1
 #endif /* _LP64 */
+#elif defined(__NetBSD__)
+#define GPR_PLATFORM_STRING "netbsd"
+#ifndef _BSD_SOURCE
+#define _BSD_SOURCE
+#endif
+#define GPR_NETBSD 1
+#define GPR_CPU_POSIX 1
+#define GPR_GCC_ATOMIC 1
+#define GPR_GCC_TLS 1
+#define GPR_POSIX_LOG 1
+#define GPR_POSIX_ENV 1
+#define GPR_POSIX_TMPFILE 1
+#define GPR_POSIX_STRING 1
+#define GPR_POSIX_SUBPROCESS 1
+#define GPR_POSIX_SYNC 1
+#define GPR_POSIX_TIME 1
+#define GPR_GETPID_IN_UNISTD_H 1
+#define GPR_SUPPORT_CHANNELS_FROM_FD 1
+#ifdef _LP64
+#define GPR_ARCH_64 1
+#else /* _LP64 */
+#define GPR_ARCH_32 1
+#endif /* _LP64 */
 #elif defined(__native_client__)
 #define GPR_PLATFORM_STRING "nacl"
 #ifndef _BSD_SOURCE
