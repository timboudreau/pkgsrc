# $NetBSD: Makefile,v 1.11 2019/04/26 14:12:47 maya Exp $

.include "../../textproc/yamcha/Makefile.common"

PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
PKGREVISION=	1
CATEGORIES+=	python

COMMENT=	YamCha python module

USE_LANGUAGES=	c c++

PYSETUPSUBDIR=	python
PYTHON_VERSIONS_INCOMPATIBLE=	36 37 # not yet ported as of 0.33

EGG_NAME=	${DISTNAME:S/-/_python-/}

BUILDLINK_API_DEPENDS.yamcha+=	yamcha>=${PKGVERSION_NOREV}
.include "../../textproc/yamcha/buildlink3.mk"
.include "../../lang/python/distutils.mk"
.include "../../mk/bsd.pkg.mk"
