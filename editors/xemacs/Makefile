# $NetBSD: Makefile,v 1.123 2019/04/11 16:05:24 hauke Exp $

PKGNAME=	${DISTNAME}
PKGREVISION=	5
COMMENT=	XEmacs text editor version 21.4

# extra options for x11 support, not for sharing with xemacs-current-nox11
# Version information in ../../editors/xemacs/Makefile.common

PKG_SUPPORTED_OPTIONS+=		x11 xface
PKG_SUGGESTED_OPTIONS+=		x11

.include "options.mk"

.include "Makefile.common"

.include "../../mk/bsd.pkg.mk"
