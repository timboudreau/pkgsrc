# $NetBSD: Makefile,v 1.191 2017/04/24 07:18:59 wiz Exp $

NO_CHECKSUM=	yes
DISTNAME=	wav2png-0.8
CATEGORIES=     graphics
X_GITHUB_TAG=    80e277e8d8ad9486b9e3f7b68f3b7d7a14a49021
MASTER_SITES=   -https://codeload.github.com/beschulz/wav2png/tar.gz/${X_GITHUB_TAG}
MAINTAINER=     tim@timboudreau.com
HOMEPAGE=       https://github.com/beschulz/wav2png/
COMMENT=        Generates soundcloud-style thumbnails from 16-bit WAV files
LICENSE=        gnu-gpl-v2
USE_TOOLS+=	gmake gtar
USE_LANGUAGES+=	c c++
BUILD_MAKE_FLAGS+=	-w
BUILD_TARGET=	all
TEST_TARGET=    check
USE_LIBTOOL=    yes
CXXFLAGS+=	-ffast-math
BUILD_DIRS=	build

post-extract:
	mv ${WRKDIR}/wav2png-${X_GITHUB_TAG} ${WRKDIR}/${DISTNAME}

do-install:
	mkdir -p ${DESTDIR}/${PREFIX}/bin
	cp ${WRKSRC}/bin/SunOS/wav2png ${DESTDIR}/${PREFIX}/bin/wav2png

#LIBS_DIRS=      -L${BUILDLINK_PREFIX.boost-libs:Q}/lib
#LIBS_DIRS+=     ${COMPILER_RPATH_FLAG}${BUILDLINK_PREFIX.boost-libs:Q}/lib

.include "../../mk/bsd.prefs.mk"
.include "../../devel/boost-headers/buildlink3.mk"
.include "../../devel/boost-libs/buildlink3.mk"
.include "../../graphics/png/buildlink3.mk"
.include "../../graphics/png++/buildlink3.mk"
.include "../../multimedia/libogg/buildlink3.mk"
.include "../../audio/libvorbis/buildlink3.mk"
.include "../../audio/libsndfile/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
