# $NetBSD: Makefile,v 1.10 2019/04/25 07:33:12 maya Exp $
#

DISTNAME=	nicovideo-dl-0.0.20190126
PKGREVISION=	1
CATEGORIES=	net
MASTER_SITES=	${MASTER_SITE_OSDN:=nicovideo-dl/70568/}

MAINTAINER=	ryoon@NetBSD.org
HOMEPAGE=	http://osdn.jp/projects/nicovideo-dl/
COMMENT=	Download videos from www.nicovideo.jp
LICENSE=	mit

USE_LANGUAGES=	# none
NO_BUILD=	yes

PYTHON_VERSIONS_INCOMPATIBLE=	27

DEPENDS+=		${PYPKGPREFIX}-expat-[0-9]*:../../textproc/py-expat

INSTALLATION_DIRS=	bin

REPLACE_PYTHON+=	nicovideo-dl

do-install:
	${INSTALL_SCRIPT} ${WRKSRC}/nicovideo-dl ${DESTDIR}${PREFIX}/bin

.include "../../lang/python/application.mk"
.include "../../mk/bsd.pkg.mk"
