$NetBSD: patch-mozilla-release_gfx_angle_checkout_src_common_third__party_smhasher_src_PMurHash.cpp,v 1.1 2019/04/19 14:02:03 fox Exp $

Fix build where _LITTLE_ENDIAN is not an integer.

Original patch from ryoon, imported from www/firefox

--- mozilla-release/gfx/angle/checkout/src/common/third_party/smhasher/src/PMurHash.cpp.orig	2018-11-16 08:40:07.000000000 +0000
+++ mozilla-release/gfx/angle/checkout/src/common/third_party/smhasher/src/PMurHash.cpp
@@ -92,7 +92,7 @@ on big endian machines, or a byte-by-byt
 /* gcc 'may' define __LITTLE_ENDIAN__ or __BIG_ENDIAN__ to 1 (Note the trailing __),
  * or even _LITTLE_ENDIAN or _BIG_ENDIAN (Note the single _ prefix) */
 #if !defined(__BYTE_ORDER)
-  #if defined(__LITTLE_ENDIAN__) && __LITTLE_ENDIAN__==1 || defined(_LITTLE_ENDIAN) && _LITTLE_ENDIAN==1
+  #if defined(__LITTLE_ENDIAN__) && __LITTLE_ENDIAN__-0==1 || defined(_LITTLE_ENDIAN) && _LITTLE_ENDIAN-0==1
     #define __BYTE_ORDER __LITTLE_ENDIAN
   #elif defined(__BIG_ENDIAN__) && __BIG_ENDIAN__==1 || defined(_BIG_ENDIAN) && _BIG_ENDIAN==1
     #define __BYTE_ORDER __BIG_ENDIAN
