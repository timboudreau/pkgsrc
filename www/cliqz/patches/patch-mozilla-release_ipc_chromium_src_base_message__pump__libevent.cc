$NetBSD: patch-mozilla-release_ipc_chromium_src_base_message__pump__libevent.cc,v 1.1 2019/04/19 14:02:03 fox Exp $

Allow older libevent

Original patch from ryoon, imported from www/firefox

--- mozilla-release/ipc/chromium/src/base/message_pump_libevent.cc.orig	2018-11-16 08:40:07.000000000 +0000
+++ mozilla-release/ipc/chromium/src/base/message_pump_libevent.cc
@@ -22,6 +22,7 @@
 
 // This macro checks that the _EVENT_SIZEOF_* constants defined in
 // ipc/chromiume/src/third_party/<platform>/event2/event-config.h are correct.
+#if 0
 #if defined(_EVENT_SIZEOF_SHORT)
 #  define CHECK_EVENT_SIZEOF(TYPE, type)                \
     static_assert(_EVENT_SIZEOF_##TYPE == sizeof(type), \
@@ -41,6 +42,7 @@ CHECK_EVENT_SIZEOF(PTHREAD_T, pthread_t)
 CHECK_EVENT_SIZEOF(SHORT, short);
 CHECK_EVENT_SIZEOF(SIZE_T, size_t);
 CHECK_EVENT_SIZEOF(VOID_P, void*);
+#endif
 
 // Lifecycle of struct event
 // Libevent uses two main data structures:
