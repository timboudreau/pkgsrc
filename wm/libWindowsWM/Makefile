# $NetBSD: Makefile,v 1.6 2018/03/07 11:57:32 wiz Exp $

DISTNAME=	libWindowsWM-1.0.1
PKGREVISION=	1
CATEGORIES=	wm x11
MASTER_SITES=	${MASTER_SITE_XORG:=lib/}
EXTRACT_SUFX=	.tar.bz2

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://xorg.freedesktop.org/
COMMENT=	Cygwin/X rootless window management extension
LICENSE=	mit

GNU_CONFIGURE=	yes
USE_LIBTOOL=	yes
USE_LANGUAGES+=	c c++
USE_TOOLS+=	gsed pkg-config

PKGCONFIG_OVERRIDE+=	windowswm.pc.in

.include "../../mk/bsd.prefs.mk"

.if !empty(USE_CROSS_COMPILE:M[yY][eE][sS])
# No harm in enabling this; the only harm might come from disabling it
# incorrectly.
CONFIGURE_ARGS+=	--enable-malloc0returnsnull
.endif

.include "../../x11/xorgproto/buildlink3.mk"
.include "../../x11/libXext/buildlink3.mk"
.include "../../x11/libX11/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
