$NetBSD: patch-src_intel_tools_aub__mem.c,v 1.1 2019/01/10 22:30:11 tnn Exp $

* Partially implement memfd_create() via mkostemp()

--- src/intel/tools/aub_mem.c.orig	2018-12-11 21:13:57.000000000 +0000
+++ src/intel/tools/aub_mem.c
@@ -34,7 +34,29 @@
 static inline int
 memfd_create(const char *name, unsigned int flags)
 {
+#if defined(__linux__)
    return syscall(SYS_memfd_create, name, flags);
+#elif defined(__FreeBSD__)
+   return shm_open(SHM_ANON, flags | O_RDWR | O_CREAT, 0600);
+#else /* DragonFly, NetBSD, OpenBSD, Solaris */
+   char template[] = "/tmp/shmfd-XXXXXX";
+#ifdef HAVE_MKOSTEMP
+   int fd = mkostemp(template, flags);
+#else
+   int fd = mkstemp(template);
+   if (flags & O_CLOEXEC) {
+      int flags = fcntl(fd, F_GETFD);
+      if (flags != -1) {
+         flags |= FD_CLOEXEC;
+         (void) fcntl(fd, F_SETFD, &flags);
+      }
+   }
+#endif /* HAVE_MKOSTEMP */
+   if (fd >= 0)
+      unlink(template);
+
+   return fd;
+#endif /* __linux__ */
 }
 #endif
 
