$NetBSD: patch-src_extension_internal_pdfinput_svg-builder.cpp,v 1.2 2019/01/06 08:41:01 markd Exp $

support for poppler 0.72 from upstream by way of linuxfromscratch 

--- src/extension/internal/pdfinput/svg-builder.cpp.orig	2018-03-11 20:38:09.000000000 +0000
+++ src/extension/internal/pdfinput/svg-builder.cpp
@@ -625,7 +625,7 @@ gchar *SvgBuilder::_createPattern(GfxPat
     if ( pattern != NULL ) {
         if ( pattern->getType() == 2 ) {  // Shading pattern
             GfxShadingPattern *shading_pattern = static_cast<GfxShadingPattern *>(pattern);
-            double *ptm;
+            const double *ptm;
             double m[6] = {1, 0, 0, 1, 0, 0};
             double det;
 
@@ -672,7 +672,7 @@ gchar *SvgBuilder::_createTilingPattern(
 
     Inkscape::XML::Node *pattern_node = _xml_doc->createElement("svg:pattern");
     // Set pattern transform matrix
-    double *p2u = tiling_pattern->getMatrix();
+    const double *p2u = tiling_pattern->getMatrix();
     double m[6] = {1, 0, 0, 1, 0, 0};
     double det;
     det = _ttm[0] * _ttm[3] - _ttm[1] * _ttm[2];    // see LP Bug 1168908
@@ -698,7 +698,7 @@ gchar *SvgBuilder::_createTilingPattern(
     pattern_node->setAttribute("patternUnits", "userSpaceOnUse");
     // Set pattern tiling
     // FIXME: don't ignore XStep and YStep
-    double *bbox = tiling_pattern->getBBox();
+    const double *bbox = tiling_pattern->getBBox();
     sp_repr_set_svg_double(pattern_node, "x", 0.0);
     sp_repr_set_svg_double(pattern_node, "y", 0.0);
     sp_repr_set_svg_double(pattern_node, "width", bbox[2] - bbox[0]);
@@ -751,7 +751,7 @@ gchar *SvgBuilder::_createTilingPattern(
  */
 gchar *SvgBuilder::_createGradient(GfxShading *shading, double *matrix, bool for_shading) {
     Inkscape::XML::Node *gradient;
-    Function *func;
+    _POPPLER_CONST Function *func;
     int num_funcs;
     bool extend0, extend1;
 
@@ -865,7 +865,7 @@ static bool svgGetShadingColorRGB(GfxSha
 
 #define INT_EPSILON 8
 bool SvgBuilder::_addGradientStops(Inkscape::XML::Node *gradient, GfxShading *shading,
-                                   Function *func) {
+                                   _POPPLER_CONST Function *func) {
     int type = func->getType();
     if ( type == 0 || type == 2 ) {  // Sampled or exponential function
         GfxRGB stop1, stop2;
@@ -877,9 +877,9 @@ bool SvgBuilder::_addGradientStops(Inksc
             _addStopToGradient(gradient, 1.0, &stop2, 1.0);
         }
     } else if ( type == 3 ) { // Stitching
-        StitchingFunction *stitchingFunc = static_cast<StitchingFunction*>(func);
-        double *bounds = stitchingFunc->getBounds();
-        double *encode = stitchingFunc->getEncode();
+        auto stitchingFunc = static_cast<_POPPLER_CONST StitchingFunction*>(func);
+        const double *bounds = stitchingFunc->getBounds();
+        const double *encode = stitchingFunc->getEncode();
         int num_funcs = stitchingFunc->getNumFuncs();
 
         // Add stops from all the stitched functions
@@ -890,7 +890,7 @@ bool SvgBuilder::_addGradientStops(Inksc
             svgGetShadingColorRGB(shading, bounds[i + 1], &color);
             // Add stops
             if (stitchingFunc->getFunc(i)->getType() == 2) {    // process exponential fxn
-                double expE = (static_cast<ExponentialFunction*>(stitchingFunc->getFunc(i)))->getE();
+                double expE = (static_cast<_POPPLER_CONST ExponentialFunction*>(stitchingFunc->getFunc(i)))->getE();
                 if (expE > 1.0) {
                     expE = (bounds[i + 1] - bounds[i])/expE;    // approximate exponential as a single straight line at x=1
                     if (encode[2*i] == 0) {    // normal sequence
@@ -1022,7 +1022,7 @@ void SvgBuilder::updateFont(GfxState *st
     if (font->getName()) {
         _font_specification = font->getName()->getCString();
     } else {
-        _font_specification = (char*) "Arial";
+        _font_specification = "Arial";
     }
 
     // Prune the font name to get the correct font family name
@@ -1030,7 +1030,7 @@ void SvgBuilder::updateFont(GfxState *st
     char *font_family = NULL;
     char *font_style = NULL;
     char *font_style_lowercase = NULL;
-    char *plus_sign = strstr(_font_specification, "+");
+    const char *plus_sign = strstr(_font_specification, "+");
     if (plus_sign) {
         font_family = g_strdup(plus_sign + 1);
         _font_specification = plus_sign + 1;
@@ -1148,7 +1148,7 @@ void SvgBuilder::updateFont(GfxState *st
     Inkscape::CSSOStringStream os_font_size;
     double css_font_size = _font_scaling * state->getFontSize();
     if ( font->getType() == fontType3 ) {
-        double *font_matrix = font->getFontMatrix();
+        const double *font_matrix = font->getFontMatrix();
         if ( font_matrix[0] != 0.0 ) {
             css_font_size *= font_matrix[3] / font_matrix[0];
         }
@@ -1193,7 +1193,7 @@ void SvgBuilder::updateTextPosition(doub
 void SvgBuilder::updateTextMatrix(GfxState *state) {
     _flushText();
     // Update text matrix
-    double *text_matrix = state->getTextMat();
+    const double *text_matrix = state->getTextMat();
     double w_scale = sqrt( text_matrix[0] * text_matrix[0] + text_matrix[2] * text_matrix[2] );
     double h_scale = sqrt( text_matrix[1] * text_matrix[1] + text_matrix[3] * text_matrix[3] );
     double max_scale;
@@ -1361,7 +1361,7 @@ void SvgBuilder::_flushText() {
     _glyphs.clear();
 }
 
-void SvgBuilder::beginString(GfxState *state, GooString * /*s*/) {
+void SvgBuilder::beginString(GfxState *state) {
     if (_need_font_update) {
         updateFont(state);
     }
