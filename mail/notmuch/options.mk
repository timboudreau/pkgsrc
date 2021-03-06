# $NetBSD: options.mk,v 1.3 2019/04/05 09:59:01 ryoon Exp $

PKG_OPTIONS_VAR=	PKG_OPTIONS.notmuch
PKG_SUPPORTED_OPTIONS=	notmuch-emacs doc
PKG_SUGGESTED_OPTIONS+=	doc

.include "../../mk/bsd.options.mk"

###
### emacs support
###
.include "../../sysutils/desktop-file-utils/desktopdb.mk"
PLIST_VARS+=		emacs
.if !empty(PKG_OPTIONS:Mnotmuch-emacs)
BUILD_DEPENDS+=		emacs>=24.0:../../editors/emacs
PLIST.emacs=		yes
.else
CONFIGURE_ARGS+=	--without-emacs
.endif

###
### documentation support
###
PLIST_VARS+=		doc
.if !empty(PKG_OPTIONS:Mdoc)
PLIST.doc=		yes
BUILD_DEPENDS+=		${PYPKGPREFIX}-sphinx-[0-9]*:../../textproc/py-sphinx
.else
CONFIGURE_ARGS+=	--without-docs
CONFIGURE_ARGS+=	--without-api-docs
.endif
