# $NetBSD: options.mk,v 1.1 2019/04/13 12:08:04 nia Exp $

PKG_OPTIONS_VAR=	PKG_OPTIONS.powder-toy
PKG_SUPPORTED_OPTIONS=	debug fftw lua
PKG_SUGGESTED_OPTIONS=	fftw lua

.include "../../mk/bsd.fast.prefs.mk"

.if ${MACHINE_ARCH} == "i386" || ${MACHINE_ARCH} == "x86_64"
PKG_SUPPORTED_OPTIONS+=	simd
.endif

.if ${MACHINE_ARCH} == "x86_64"
PKG_SUGGESTED_OPTIONS+=	simd
.endif

.include "../../mk/bsd.options.mk"

.if !empty(PKG_OPTIONS:Mdebug)
SCONS_ARGS+=	--debugging
.else
SCONS_ARGS+=	--release
.endif

.if !empty(PKG_OPTIONS:Mfftw)
.include "../../math/fftwf/buildlink3.mk"
.else
SCONS_ARGS+=	--nofft
.endif

.if !empty(PKG_OPTIONS:Mlua)
.include "../../lang/lua51/buildlink3.mk"
.else
SCONS_ARGS+=	--nolua
.endif

.if !empty(PKG_OPTIONS:Msimd)
SCONS_ARGS+=	--sse
SCONS_ARGS+=	--sse2
.else
SCONS_ARGS+=	--no-sse
.endif
