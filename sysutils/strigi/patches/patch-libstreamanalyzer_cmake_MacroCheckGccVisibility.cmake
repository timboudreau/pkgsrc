$NetBSD: patch-libstreamanalyzer_cmake_MacroCheckGccVisibility.cmake,v 1.2 2018/02/08 07:28:58 mrg Exp $

Support gcc6/7/8/9

--- libstreamanalyzer/cmake/MacroCheckGccVisibility.cmake.orig	2013-02-05 21:34:52.000000000 +0000
+++ libstreamanalyzer/cmake/MacroCheckGccVisibility.cmake
@@ -15,7 +15,7 @@ macro(MACRO_CHECK_GCC_VISIBILITY GccVisi
    # get the gcc version
    exec_program(${CMAKE_C_COMPILER} ARGS ${CMAKE_C_COMPILER_ARG1} --version OUTPUT_VARIABLE _gcc_version_info)
 
-   string (REGEX MATCH "[345]\\.[0-9]\\.[0-9]" _gcc_version "${_gcc_version_info}")
+   string (REGEX MATCH "[3456789]\\.[0-9]\\.[0-9]" _gcc_version "${_gcc_version_info}")
    # gcc on mac just reports: "gcc (GCC) 3.3 20030304 ..." without the patch level, handle this here:
    if (NOT _gcc_version)
       string (REGEX REPLACE ".*\\(GCC\\).* ([34]\\.[0-9]) .*" "\\1.0" _gcc_version "${_gcc_version_info}")
