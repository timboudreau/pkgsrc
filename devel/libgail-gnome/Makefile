# $NetBSD: Makefile,v 1.97 2019/04/03 00:32:37 ryoon Exp $

DISTNAME=	libgail-gnome-1.20.4
PKGREVISION=	43
CATEGORIES=	devel gnome
MASTER_SITES=	${MASTER_SITE_GNOME:=sources/libgail-gnome/1.20/}
EXTRACT_SUFX=	.tar.bz2

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://developer.gnome.org/accessibility-devel-guide/stable/dev-start-5.html.en
COMMENT=	GNOME Accessibility Implementation Library

GNU_CONFIGURE=		YES
USE_LIBTOOL=		YES
USE_TOOLS+=		pkg-config gmake

PKGCONFIG_OVERRIDE=	libgail-gnome.pc.in
GCONF_SCHEMAS=		libgail-gnome.schemas

.include "../../devel/GConf/schemas.mk"
.include "../../devel/at-spi/buildlink3.mk"
.include "../../devel/atk/buildlink3.mk"
.include "../../devel/libbonobo/buildlink3.mk"
.include "../../devel/libbonoboui/buildlink3.mk"
.include "../../devel/libgnomeui/buildlink3.mk"
.include "../../x11/gnome-panel/buildlink3.mk"
.include "../../x11/gtk2/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
