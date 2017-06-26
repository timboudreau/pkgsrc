$NetBSD$

Disable __stack_chk_fail_local 32-bit optimisation on SunOS.

--- gcc/targhooks.c.orig	2014-03-03 21:51:58.000000000 +0000
+++ gcc/targhooks.c
@@ -716,7 +716,7 @@ default_external_stack_protect_fail (voi
 tree
 default_hidden_stack_protect_fail (void)
 {
-#ifndef HAVE_GAS_HIDDEN
+#if !defined(HAVE_GAS_HIDDEN) || defined(__sun)
   return default_external_stack_protect_fail ();
 #else
   tree t = stack_chk_fail_decl;
