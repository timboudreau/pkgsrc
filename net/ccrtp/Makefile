# $NetBSD: Makefile,v 1.13 2019/01/13 18:22:10 adam Exp $

DISTNAME=	ccrtp-2.1.2
PKGREVISION=	5
CATEGORIES=	net
MASTER_SITES=	${MASTER_SITE_GNU:=ccrtp/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://www.gnu.org/software/ccrtp/
COMMENT=	RTP and RTSP protocol implementation using GNU CommonCpp
LICENSE=	gnu-gpl-v2

USE_CMAKE=		yes
USE_LANGUAGES+=		c c++
USE_PKGLOCALEDIR=	yes
USE_TOOLS+=		pkg-config
PKGCONFIG_OVERRIDE+=	libccrtp.pc.in
INFO_FILES=		yes

CMAKE_ARGS+=	-DLIB_VERSION=${PKGVERSION_NOREV}
CMAKE_ARGS+=	-DSOVERSION=2

.include "options.mk"

BUILDLINK_API_DEPENDS.ucommon+=	ucommon>=6.3.1
.include "../../devel/ucommon/buildlink3.mk"
.include "../../mk/pthread.buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
