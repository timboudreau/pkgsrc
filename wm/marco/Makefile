# $NetBSD: Makefile,v 1.16 2019/04/25 07:33:26 maya Exp $

.include "../../meta-pkgs/mate/Makefile.common"

DISTNAME=	marco-${VERSION:R}.0
PKGREVISION=	1
CATEGORIES=	wm

COMMENT=	Fork of GNOME Metacity

USE_LANGUAGES=	c99

BUILD_DEPENDS+=		gdk-pixbuf2-[0-9]*:../../graphics/gdk-pixbuf2
TOOL_DEPENDS+=		itstool-[0-9]*:../../textproc/itstool
TOOL_DEPENDS+=		libxml2-[0-9]*:../../textproc/libxml2 # for xmllint
DEPENDS+=		zenity-[0-9]*:../../x11/zenity

PKGCONFIG_OVERRIDE+=	src/libmarco-private.pc.in

.include "../../sysutils/desktop-file-utils/desktopdb.mk"
.include "../../audio/libcanberra/buildlink3.mk"
.include "../../sysutils/libgtop/buildlink3.mk"
.include "../../x11/startup-notification/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
