# $NetBSD: Makefile,v 1.85 2019/04/03 00:32:55 ryoon Exp $

DISTNAME=	fuppes-0.660
PKGREVISION=	49
CATEGORIES=	multimedia
MASTER_SITES=	${MASTER_SITE_SOURCEFORGE:=fuppes/}

MAINTAINER=	tech-multimedia@NetBSD.org
HOMEPAGE=	http://fuppes.ulrich-voelkel.de/
COMMENT=	Free UPnP Entertainment Service
LICENSE=	gnu-gpl-v2

MAKE_JOBS_SAFE=	no

GNU_CONFIGURE=	yes
USE_LANGUAGES=	c c++03
USE_TOOLS+=	gmake pkg-config
USE_LIBTOOL=	yes

CONFIGURE_ARGS+=	--enable-faad
CONFIGURE_ARGS+=	--enable-mad
CONFIGURE_ARGS+=	--enable-lame

.include "options.mk"

.include "../../devel/pcre/buildlink3.mk"
.include "../../audio/faad2/buildlink3.mk"
.include "../../audio/flac/buildlink3.mk"
.include "../../audio/lame/buildlink3.mk"
.include "../../audio/libmad/buildlink3.mk"
.include "../../audio/libmpcdec/buildlink3.mk"
.include "../../audio/libvorbis/buildlink3.mk"
.include "../../audio/taglib/buildlink3.mk"
.include "../../audio/twolame/buildlink3.mk"
.include "../../databases/sqlite3/buildlink3.mk"
.include "../../graphics/ImageMagick6/buildlink3.mk"
.include "../../graphics/exiv2/buildlink3.mk"
.include "../../textproc/libxml2/buildlink3.mk"
.include "../../graphics/simage/buildlink3.mk"
.include "../../multimedia/ffmpeg1/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
