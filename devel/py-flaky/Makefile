# $NetBSD: Makefile,v 1.2 2019/01/21 09:34:53 adam Exp $

DISTNAME=	flaky-3.5.3
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	devel python
MASTER_SITES=	${MASTER_SITE_PYPI:=f/flaky/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/box/flaky
COMMENT=	Plugin for nose or py.test that automatically reruns flaky tests
LICENSE=	apache-2.0

TEST_DEPENDS+=	${PYPKGPREFIX}-tox-[0-9]*:../../devel/py-tox

USE_LANGUAGES=	# none

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
