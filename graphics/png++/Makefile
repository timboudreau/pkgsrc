# $NetBSD: Makefile,v 1.191 2017/04/24 07:18:59 wiz Exp $

DISTNAME=       png++-0.2.9
CATEGORIES=     graphics
MASTER_SITES=	http://download-mirror.savannah.gnu.org/releases/pngpp/
MASTER_SITES+=	http://www.namesdir.com/mirrors/nongnu/pngpp/
MASTER_SITES+=	http://gnu.mirrors.pair.com/savannah/savannah/pngpp/
MASTER_SITES+=	http://savannah.spinellicreations.com/pngpp/
MAINTAINER=     tim@timboudreau.com
HOMEPAGE=       http://savannah.nongnu.org/projects/pngpp
COMMENT=        C++ wrapper for libpng
LICENSE=        2-clause-bsd
USE_TOOLS+=	gmake
USE_LANGUAGES+=	c c++
PATCH_DIST_STRIP=-p0

BUILD_TARGET=	install
PKGCONFIG_OVERRIDE=
TEST_TARGET=    check
USE_LIBTOOL=    yes

.if ${BUILD_TARGET} == "do-install"
MAKEFLAGS+=	-w PREFIX=../.destdir/opt/local
.endif

CPPFLAGS+=      -DPNG_USER_CHUNK_CACHE_MAX=128 -DPNG_USER_CHUNK_MALLOC_MAX=8000000

.include "../../mk/bsd.prefs.mk"
.include "../../graphics/png/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
