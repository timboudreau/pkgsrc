$NetBSD: patch-toolkit_xre_glxtest.cpp,v 1.2 2019/02/25 15:32:24 wiz Exp $

Fix libGL filename on NetBSD,
see https://bugzilla.mozilla.org/show_bug.cgi?id=1180498

--- toolkit/xre/glxtest.cpp.orig	2019-02-23 20:00:59.000000000 +0000
+++ toolkit/xre/glxtest.cpp
@@ -116,7 +116,7 @@ void glxtest() {
         "The MOZ_AVOID_OPENGL_ALTOGETHER environment variable is defined");
 
     ///// Open libGL and load needed symbols /////
-#ifdef __OpenBSD__
+#if defined(__OpenBSD__) || defined(__NetBSD__)
 #define LIBGL_FILENAME "libGL.so"
 #else
 #define LIBGL_FILENAME "libGL.so.1"
