# $NetBSD: Makefile,v 1.12 2019/04/26 14:12:44 maya Exp $

DISTNAME=	socketpool-0.5.3
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	net
MASTER_SITES=	${MASTER_SITE_PYPI:=s/socketpool/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/benoitc/socketpool
COMMENT=	Simple Python socket pool
LICENSE=	unlicense OR mit

USE_LANGUAGES=	# none

DOCDIR=			share/doc/${PKGBASE:S/${PYPKGPREFIX}/py/1}
INSTALLATION_DIRS+=	${DOCDIR}

PYTHON_VERSIONS_INCOMPATIBLE=	36 37 # not yet ported as of 0.5.3

post-install:
	${INSTALL_DATA} ${DESTDIR}${PREFIX}/socketpool/README.rst \
		${DESTDIR}${PREFIX}/${DOCDIR}
	${RM} -rf ${DESTDIR}${PREFIX}/socketpool

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
