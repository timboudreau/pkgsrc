$NetBSD: patch-servo_components_style_lib.rs,v 1.1 2019/03/15 11:51:26 wiz Exp $

--- servo/components/style/lib.rs.orig	2019-02-21 19:22:50.000000000 +0000
+++ servo/components/style/lib.rs
@@ -23,8 +23,6 @@
 //! [cssparser]: ../cssparser/index.html
 //! [selectors]: ../selectors/index.html
 
-#![deny(missing_docs)]
-
 extern crate app_units;
 extern crate arrayvec;
 extern crate atomic_refcell;
