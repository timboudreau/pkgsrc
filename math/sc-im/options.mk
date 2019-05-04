# $NetBSD: options.mk,v 1.2 2019/02/11 22:01:18 leot Exp $

PKG_OPTIONS_VAR=	PKG_OPTIONS.sc-im
PKG_SUPPORTED_OPTIONS+=	color xls xlsx lua
PKG_SUGGESTED_OPTIONS+=	color xls xlsx lua

.include "../../mk/bsd.options.mk"

.if !empty(PKG_OPTIONS:Mcolor)
CFLAGS+=	-DUSECOLORS
.endif

.if !empty(PKG_OPTIONS:Mxls)
CFLAGS+=	-DXLS
LDFLAGS+=	-lxlsreader
.include "../../textproc/libxls/buildlink3.mk"
.endif

.if !empty(PKG_OPTIONS:Mxlsx)
CFLAGS+=	-DXLSX -DXLSX_EXPORT
LDFLAGS+=	-lxlsxwriter
.include "../../archivers/libzip/buildlink3.mk"
.include "../../textproc/libxml2/buildlink3.mk"
.include "../../textproc/libxlsxwriter/buildlink3.mk"
.endif

.if !empty(PKG_OPTIONS:Mlua)
.include "../../lang/lua51/buildlink3.mk"
.endif
